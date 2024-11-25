package telegram

import (
	"context"
	"fmt"
	"main/pkg/data"
	mutes "main/pkg/mutes"
	"main/pkg/report/entry"
	"main/pkg/state"
	"main/pkg/templates"
	"strings"
	"time"

	"github.com/guregu/null/v5"

	"go.opentelemetry.io/otel/trace"

	"main/pkg/types"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type Reporter struct {
	TelegramToken    string
	TelegramChat     int64
	MutesManager     *mutes.Manager
	StateGenerator   *state.Generator
	DataManager      *data.Manager
	TemplatesManager templates.Manager
	Tracer           trace.Tracer

	TelegramBot *tele.Bot
	Logger      zerolog.Logger

	Version string
}

const (
	MaxMessageSize = 4096
)

func NewTelegramReporter(
	config types.TelegramConfig,
	mutesManager *mutes.Manager,
	stateGenerator *state.Generator,
	dataManager *data.Manager,
	logger *zerolog.Logger,
	version string,
	timezone *time.Location,
	tracer trace.Tracer,
) *Reporter {
	return &Reporter{
		TelegramToken:    config.TelegramToken,
		TelegramChat:     config.TelegramChat,
		MutesManager:     mutesManager,
		StateGenerator:   stateGenerator,
		DataManager:      dataManager,
		Logger:           logger.With().Str("component", "telegram_reporter").Logger(),
		TemplatesManager: templates.NewTelegramTemplatesManager(logger, timezone),
		Version:          version,
		Tracer:           tracer,
	}
}

func (reporter *Reporter) Init() error {
	if reporter.TelegramToken == "" || reporter.TelegramChat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return nil
	}

	if err := reporter.InitBot(); err != nil {
		return err
	}

	go reporter.TelegramBot.Start()

	return nil
}

func (reporter *Reporter) InitBot() error {
	bot, err := tele.NewBot(tele.Settings{
		Token:  reporter.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return err
	}

	bot.Handle("/start", reporter.HandleHelp)
	bot.Handle("/help", reporter.HandleHelp)
	bot.Handle("/proposals_mute", reporter.HandleAddMute)
	bot.Handle("/proposals_unmute", reporter.HandleDeleteMute)
	bot.Handle("/proposals_mutes", reporter.HandleListMutes)
	bot.Handle("/proposals", reporter.HandleProposals)
	bot.Handle("/tally", reporter.HandleTally)
	bot.Handle("/params", reporter.HandleParams)

	reporter.TelegramBot = bot

	return nil
}

func (reporter *Reporter) Enabled() bool {
	return reporter.TelegramToken != "" && reporter.TelegramChat != 0
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

	_, err = reporter.TelegramBot.Send(
		&tele.User{
			ID: reporter.TelegramChat,
		},
		serializedEntry,
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		reporter.Logger.Err(err).Msg("Could not send Telegram message")
		return err
	}

	return nil
}

func (reporter *Reporter) Name() string {
	return "telegram-reporter"
}

func (reporter *Reporter) BotReply(c tele.Context, msg string) error {
	msgsByNewline := strings.Split(msg, "\n")

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > MaxMessageSize {
			if err := c.Reply(sb.String(), tele.ModeHTML, tele.NoPreview); err != nil {
				reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
				return err
			}

			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	if err := c.Reply(sb.String(), tele.ModeHTML, tele.NoPreview); err != nil {
		reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
		return err
	}

	return nil
}

func ParseMuteOptions(query string, c tele.Context) (*types.Mute, string) {
	args := strings.Split(query, " ")
	if len(args) < 2 {
		return nil, "Usage: /proposals_mute <duration> [params]"
	}

	_, durationString, args := args[0], args[1], args[2:] // removing first argument as it's always /proposals_mute

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return nil, fmt.Sprintf("Invalid duration provided: %s", durationString)
	}

	mute := &types.Mute{
		Chain:      null.NewString("", false),
		ProposalID: null.NewString("", false),
		Expires:    time.Now().Add(duration),
		Comment: fmt.Sprintf(
			"Muted using cosmos-proposals-checker for %s by %s",
			duration,
			c.Sender().FirstName,
		),
	}

	for index, arg := range args {
		argSplit := strings.SplitN(arg, "=", 2)
		if len(argSplit) < 2 {
			return nil, fmt.Sprintf(
				"Invalid param at position %d: expected an expression like \"[chain=cosmos]\", but got %s",
				index+1,
				arg,
			)
		}

		switch argSplit[0] {
		case "chain":
			mute.Chain = null.StringFrom(argSplit[1])
		case "proposal":
			mute.ProposalID = null.StringFrom(argSplit[1])
		}
	}

	return mute, ""
}

func ParseMuteDeleteOptions(c tele.Context) (*types.Mute, string) {
	// we only construct mute with chain/proposal to compare, no need to take care
	// about the expiration/comment
	mute := &types.Mute{
		Chain:      null.NewString("", false),
		ProposalID: null.NewString("", false),
		Expires:    time.Now(),
		Comment:    "",
	}

	for index, arg := range c.Args() {
		argSplit := strings.SplitN(arg, "=", 2)
		if len(argSplit) < 2 {
			return nil, fmt.Sprintf(
				"Invalid param at position %d: expected an expression like \"[chain=cosmos]\", but got %s",
				index+1,
				arg,
			)
		}

		switch argSplit[0] {
		case "chain":
			mute.Chain = null.StringFrom(argSplit[1])
		case "proposal":
			mute.ProposalID = null.StringFrom(argSplit[1])
		}
	}

	return mute, ""
}

func (reporter *Reporter) Stop() {
	reporter.Logger.Info().Msg("Shutting down...")
	reporter.TelegramBot.Stop()
}
