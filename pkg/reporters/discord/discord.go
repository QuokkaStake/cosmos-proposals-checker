package discord

import (
	"context"
	"main/pkg/data"
	mutes "main/pkg/mutes"
	"main/pkg/report/entry"
	statePkg "main/pkg/state"
	templatesPkg "main/pkg/templates"
	types "main/pkg/types"
	"main/pkg/utils"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type Reporter struct {
	Token   string
	Guild   string
	Channel string

	Version string

	DiscordSession   *discordgo.Session
	StateGenerator   *statePkg.Generator
	Logger           zerolog.Logger
	Config           *types.Config
	MutesManager     *mutes.Manager
	DataManager      *data.Manager
	TemplatesManager templatesPkg.Manager
	Commands         map[string]*Command
	Tracer           trace.Tracer
	Timezone         *time.Location
}

func NewReporter(
	config *types.Config,
	version string,
	logger *zerolog.Logger,
	mutesManager *mutes.Manager,
	dataManager *data.Manager,
	stateGenerator *statePkg.Generator,
	timezone *time.Location,
	tracer trace.Tracer,
) *Reporter {
	return &Reporter{
		Token:            config.DiscordConfig.Token,
		Guild:            config.DiscordConfig.Guild,
		Channel:          config.DiscordConfig.Channel,
		Config:           config,
		Logger:           logger.With().Str("component", "discord_reporter").Logger(),
		MutesManager:     mutesManager,
		DataManager:      dataManager,
		StateGenerator:   stateGenerator,
		TemplatesManager: templatesPkg.NewDiscordTemplatesManager(logger, timezone),
		Commands:         make(map[string]*Command, 0),
		Version:          version,
		Timezone:         timezone,
		Tracer:           tracer,
	}
}

func (reporter *Reporter) Init() error {
	if !reporter.Enabled() {
		reporter.Logger.Debug().Msg("Discord credentials not set, not creating Discord reporter")
		return nil
	}
	session, err := discordgo.New("Bot " + reporter.Token)
	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Error initializing Discord bot")
		return err
	}

	reporter.DiscordSession = session

	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Error opening Discord websocket session")
		return nil
	}

	reporter.Logger.Info().Err(err).Msg("Discord bot listening")

	reporter.Commands = map[string]*Command{
		"help":             reporter.GetHelpCommand(),
		"proposals":        reporter.GetProposalsCommand(),
		"proposals_mute":   reporter.GetAddMuteCommand(),
		"proposals_unmute": reporter.GetDeleteMuteCommand(),
		"proposals_mutes":  reporter.GetMutesCommand(),
		"params":           reporter.GetParamsCommand(),
		"tally":            reporter.GetTallyCommand(),
	}

	go reporter.InitCommands()
	return nil
}

func (reporter *Reporter) InitCommands() {
	session := reporter.DiscordSession
	var wg sync.WaitGroup
	var mutex sync.Mutex

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name

		if command, ok := reporter.Commands[commandName]; ok {
			command.Handler(s, i)
		}
	})

	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, reporter.Guild)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Could not fetch registered commands")
		return
	}

	desiredCommands := utils.Map(
		utils.MapToArray(reporter.Commands),
		func(c *Command) *discordgo.ApplicationCommand { return c.Info },
	)

	commandsToAdd := utils.Subtract(desiredCommands, registeredCommands, func(v *discordgo.ApplicationCommand) string {
		return v.Name
	})

	commandsToDelete := utils.Subtract(registeredCommands, desiredCommands, func(v *discordgo.ApplicationCommand) string {
		return v.Name
	})

	commandsToUpdate := utils.Union(registeredCommands, desiredCommands, func(v *discordgo.ApplicationCommand) string {
		return v.Name
	})

	reporter.Logger.Info().
		Int("commands_to_add", len(commandsToAdd)).
		Int("commands_to_delete", len(commandsToDelete)).
		Int("commands_to_update", len(commandsToUpdate)).
		Msg("Updating Discord slash commands")

	wg.Add(len(commandsToAdd) + len(commandsToDelete) + len(commandsToUpdate))

	for _, command := range commandsToDelete {
		go func(command *discordgo.ApplicationCommand) {
			defer wg.Done()

			err := session.ApplicationCommandDelete(session.State.User.ID, reporter.Guild, command.ID)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("command", command.Name).Msg("Could not delete command")
				return
			}
			reporter.Logger.Info().Str("command", command.Name).Msg("Deleted command")
		}(command)
	}

	for _, command := range commandsToAdd {
		go func(command *discordgo.ApplicationCommand) {
			defer wg.Done()

			cmd, err := session.ApplicationCommandCreate(session.State.User.ID, reporter.Guild, command)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("command", command.Name).Msg("Could not create command")
				return
			}
			reporter.Logger.Info().Str("command", cmd.Name).Msg("Created command")

			mutex.Lock()
			reporter.Commands[command.Name].Info = cmd
			mutex.Unlock()
		}(command)
	}

	for _, command := range commandsToUpdate {
		go func(command *discordgo.ApplicationCommand) {
			defer wg.Done()

			cmd, err := session.ApplicationCommandEdit(
				session.State.User.ID,
				reporter.Guild,
				command.ID,
				command,
			)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("command", command.Name).Msg("Could not update command")
				return
			}
			reporter.Logger.Info().Str("command", cmd.Name).Msg("Updated command")

			mutex.Lock()
			reporter.Commands[command.Name].Info = cmd
			mutex.Unlock()
		}(command)
	}

	wg.Wait()
	reporter.Logger.Info().Msg("All commands updated")
}

func (reporter *Reporter) Enabled() bool {
	return reporter.Token != "" && reporter.Guild != "" && reporter.Channel != ""
}

func (reporter *Reporter) Name() string {
	return "discord-reporter"
}

func (reporter *Reporter) SerializeReportEntry(e entry.ReportEntry) (string, error) {
	return reporter.TemplatesManager.Render(e.Name(), e)
}

func (reporter *Reporter) SendReportEntry(reportEntry entry.ReportEntry, ctx context.Context) error {
	_, span := reporter.Tracer.Start(ctx, "Sending Telegram report entry")
	defer span.End()

	serializedEntry, err := reporter.SerializeReportEntry(reportEntry)
	if err != nil {
		reporter.Logger.Err(err).Msg("Could not serialize report entry")
		return err
	}

	_, err = reporter.DiscordSession.ChannelMessageSend(
		reporter.Channel,
		serializedEntry,
	)

	return err
}

func (reporter *Reporter) BotRespond(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	chunks := utils.SplitStringIntoChunks(text, 2000)
	firstChunk, rest := chunks[0], chunks[1:]

	reporter.BotSendInteraction(s, i, firstChunk)

	for _, chunk := range rest {
		reporter.BotSendFollowup(s, i, chunk)
	}
}

func (reporter *Reporter) BotSendInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	}); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error sending response")
	}
}

func (reporter *Reporter) BotSendFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: text,
	}); err != nil {
		reporter.Logger.Error().
			Err(err).
			Msg("Error sending followup message")
	}
}

func (reporter *Reporter) SerializeDate(date time.Time) string {
	return date.Format(time.RFC822)
}
