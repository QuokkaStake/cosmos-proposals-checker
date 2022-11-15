package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type TelegramReporter struct {
	TelegramToken  string
	TelegramChat   int64
	MutesManager   *MutesManager
	StateGenerator *StateGenerator

	TelegramBot *tele.Bot
	Logger      zerolog.Logger
	Templates   map[ReportEntryType]*template.Template
}

const (
	MaxMessageSize = 4096
)

//go:embed templates/*
var templatesFs embed.FS

func NewTelegramReporter(
	config TelegramConfig,
	mutesManager *MutesManager,
	stateGenerator *StateGenerator,
	logger *zerolog.Logger) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken:  config.TelegramToken,
		TelegramChat:   config.TelegramChat,
		MutesManager:   mutesManager,
		StateGenerator: stateGenerator,
		Logger:         logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:      make(map[ReportEntryType]*template.Template, 0),
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

func (reporter TelegramReporter) Enabled() bool {
	return reporter.TelegramToken != "" && reporter.TelegramChat != 0
}

func (reporter TelegramReporter) GetTemplate(t ReportEntryType) (*template.Template, error) {
	if template, ok := reporter.Templates[t]; ok {
		reporter.Logger.Trace().Str("type", string(t)).Msg("Using cached template")
		return template, nil
	}

	reporter.Logger.Trace().Str("type", string(t)).Msg("Loading template")

	filename := fmt.Sprintf("templates/telegram/%s.html", t)
	template, err := template.ParseFS(templatesFs, filename)
	if err != nil {
		return nil, err
	}

	reporter.Templates[t] = template

	return template, nil
}

func (reporter *TelegramReporter) SerializeReportEntry(e ReportEntry) (string, error) {
	template, err := reporter.GetTemplate(e.Type)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", string(e.Type)).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = template.Execute(&buffer, e)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", string(e.Type)).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), nil
}

func (reporter TelegramReporter) SendReport(report Report) error {
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

func (reporter TelegramReporter) Name() string {
	return "telegram-reporter"
}

func (reporter *TelegramReporter) HandleAddMute(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got add mute query")

	mute, err := ParseMuteOptions(c.Text(), c)
	if err != "" {
		return c.Reply(fmt.Sprintf("Error muting notification: %s", err))
	}

	reporter.MutesManager.AddMute(mute)
	if mute.ProposalID != "" {
		return reporter.BotReply(c, fmt.Sprintf(
			"Notification for proposal #%s on %s are muted till %s.",
			mute.ProposalID,
			mute.Chain,
			mute.GetExpirationTime(),
		))
	}

	return reporter.BotReply(c, fmt.Sprintf(
		"Notification for all proposals on %s are muted till %s.",
		mute.Chain,
		mute.GetExpirationTime(),
	))
}

func (reporter *TelegramReporter) HandleListMutes(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list mutes query")

	mutes := Filter(reporter.MutesManager.Mutes.Mutes, func(m Mute) bool {
		return !m.IsExpired()
	})

	template, _ := reporter.GetTemplate("mutes")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, mutes); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering votes template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}

func (reporter *TelegramReporter) HandleProposals(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got proposals list query")

	state := reporter.StateGenerator.GetState(NewState())
	template, _ := reporter.GetTemplate("proposals")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, state); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering votes template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}

func (reporter *TelegramReporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	template, _ := reporter.GetTemplate("help")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, nil); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering telp template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
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

func ParseMuteOptions(query string, c tele.Context) (Mute, string) {
	args := strings.Split(query, " ")
	if len(args) <= 2 {
		return Mute{}, "Usage: /proposals_mute <duration> <chain> [<proposal>]"
	}

	_, args = args[0], args[1:] // removing first argument as it's always /proposals_mute
	durationString, chain := args[0], args[1]
	proposalID := ""
	if len(args) >= 3 {
		proposalID = args[2]
	}

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return Mute{}, "Invalid duration provided"
	}

	mute := Mute{
		Chain:      chain,
		ProposalID: proposalID,
		Expires:    time.Now().Add(duration),
		Comment: fmt.Sprintf(
			"Muted using cosmos-proposals-checker for %s by %s",
			duration,
			c.Sender().FirstName,
		),
	}

	return mute, ""
}
