package telegram

import (
	"bytes"
	"fmt"
	"html/template"
	mutes "main/pkg/mutes"
	"main/pkg/state"
	"strings"
	"time"

	"main/pkg/config"
	"main/pkg/reporters"
	"main/pkg/types"
	"main/templates"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type TelegramReporter struct {
	TelegramToken  string
	TelegramChat   int64
	MutesManager   *mutes.Manager
	StateGenerator *state.Generator

	TelegramBot *tele.Bot
	Logger      zerolog.Logger
	Templates   map[types.ReportEntryType]*template.Template
}

const (
	MaxMessageSize = 4096
)

func NewTelegramReporter(
	config config.TelegramConfig,
	mutesManager *mutes.Manager,
	stateGenerator *state.Generator,
	logger *zerolog.Logger,
) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken:  config.TelegramToken,
		TelegramChat:   config.TelegramChat,
		MutesManager:   mutesManager,
		StateGenerator: stateGenerator,
		Logger:         logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:      make(map[types.ReportEntryType]*template.Template, 0),
	}
}

func (reporter *TelegramReporter) Init() {
	if reporter.TelegramToken == "" || reporter.TelegramChat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  reporter.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return
	}

	bot.Handle("/start", reporter.HandleHelp)
	bot.Handle("/help", reporter.HandleHelp)
	bot.Handle("/proposals_mute", reporter.HandleAddMute)
	bot.Handle("/proposals_mutes", reporter.HandleListMutes)
	bot.Handle("/proposals", reporter.HandleProposals)

	reporter.TelegramBot = bot
	go reporter.TelegramBot.Start()
}

func (reporter *TelegramReporter) Enabled() bool {
	return reporter.TelegramToken != "" && reporter.TelegramChat != 0
}

func (reporter *TelegramReporter) GetTemplate(tmlpType types.ReportEntryType) (*template.Template, error) {
	if cachedTemplate, ok := reporter.Templates[tmlpType]; ok {
		reporter.Logger.Trace().Str("type", string(tmlpType)).Msg("Using cached template")
		return cachedTemplate, nil
	}

	reporter.Logger.Trace().Str("type", string(tmlpType)).Msg("Loading template")

	filename := fmt.Sprintf("%s.html", tmlpType)

	t, err := template.New(filename).Funcs(template.FuncMap{
		"SerializeLink": reporter.SerializeLink,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	reporter.Templates[tmlpType] = t

	return t, nil
}

func (reporter *TelegramReporter) SerializeReportEntry(e reporters.ReportEntry) (string, error) {
	parsedTemplate, err := reporter.GetTemplate(e.Type)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", string(e.Type)).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = parsedTemplate.Execute(&buffer, e)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", string(e.Type)).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), nil
}

func (reporter *TelegramReporter) SendReport(report reporters.Report) error {
	for _, entry := range report.Entries {
		if !entry.HasVoted() && reporter.MutesManager.IsMuted(entry.Chain.Name, entry.ProposalID) {
			reporter.Logger.Debug().
				Str("chain", entry.Chain.Name).
				Str("proposal", entry.ProposalID).
				Msg("Notifications are muted, not sending.")
			continue
		}

		serializedEntry, err := reporter.SerializeReportEntry(entry)
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
	}

	return nil
}

func (reporter *TelegramReporter) Name() string {
	return "telegram-reporter"
}

func (reporter *TelegramReporter) BotReply(c tele.Context, msg string) error {
	msgsByNewline := strings.Split(msg, "\n")

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > MaxMessageSize {
			if err := c.Reply(sb.String(), tele.ModeHTML); err != nil {
				reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
				return err
			}

			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	if err := c.Reply(sb.String(), tele.ModeHTML); err != nil {
		reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
		return err
	}

	return nil
}

func ParseMuteOptions(query string, c tele.Context) (*mutes.Mute, string) {
	args := strings.Split(query, " ")
	if len(args) < 2 {
		return nil, "Usage: /proposals_mute <duration> [params]"
	}

	_, durationString, args := args[0], args[1], args[2:] // removing first argument as it's always /proposals_mute

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return nil, fmt.Sprintf("Invalid duration provided: %s", durationString)
	}

	mute := &mutes.Mute{
		Chain:      "",
		ProposalID: "",
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
			mute.Chain = argSplit[1]
		case "proposal":
			mute.ProposalID = argSplit[1]
		}
	}

	return mute, ""
}

func (reporter *TelegramReporter) SerializeLink(link types.Link) template.HTML {
	if link.Href != "" {
		return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, link.Name))
	}

	return template.HTML(link.Name)
}
