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
	TelegramToken string
	TelegramChat  int64
	MutesManager  *MutesManager

	TelegramBot *tele.Bot
	Logger      zerolog.Logger
	Templates   map[ReportEntryType]*template.Template
}

const (
	MaxMessageSize = 4096
)

//go:embed templates/*
var templatesFs embed.FS

func NewTelegramReporter(config TelegramConfig, mutesManager *MutesManager, logger *zerolog.Logger) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken: config.TelegramToken,
		TelegramChat:  config.TelegramChat,
		MutesManager:  mutesManager,
		Logger:        logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:     make(map[ReportEntryType]*template.Template, 0),
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

// func (reporter *TelegramReporter) SerializeReportEntry(e ReportEntry) string {
// 	if e.Type == ProposalQueryError {
// 		return reporter.SerializeProposalsError(e)
// 	}
// 	if e.Type == VoteQueryError {
// 		return reporter.SerializeVoteError(e)
// 	}

// 	var sb strings.Builder

// 	messageText := "🔴 <strong>Wallet %s hasn't voted on proposal %s on %s</strong>\n%s\n\n"
// 	if e.HasRevoted() {
// 		messageText = "↔️ <strong>Wallet %s hasn changed its vote on proposal %s on %s</strong>\n%s\n\n"
// 	} else if e.HasVoted() {
// 		messageText = "✅ <strong>Wallet %s has voted on proposal %s on %s</strong>\n%s\n\n"
// 	}

// 	sb.WriteString(fmt.Sprintf(
// 		messageText,
// 		e.Wallet,
// 		e.ProposalID,
// 		e.Chain.GetName(),
// 		e.ProposalTitle,
// 	))

// 	if e.HasVoted() {
// 		sb.WriteString(fmt.Sprintf(
// 			"<strong>Vote: </strong>%s\n",
// 			e.Value,
// 		))
// 	}
// 	if e.HasRevoted() {
// 		sb.WriteString(fmt.Sprintf(
// 			"<strong>Old vote: </strong>%s\n",
// 			e.OldValue,
// 		))
// 	}

// 	sb.WriteString(fmt.Sprintf(
// 		"Voting ends at: %s (in %s)\n\n",
// 		e.ProposalVoteEndingTime.Format(time.RFC3339Nano),
// 		time.Until(e.ProposalVoteEndingTime).Round(time.Second),
// 	))

// 	sb.WriteString(reporter.SerializeLinks(e))
// 	sb.WriteString(AuthorDisclaimer)

// 	return sb.String()
// }

func (reporter TelegramReporter) SerializeLinks(e ReportEntry) string {
	var sb strings.Builder

	if e.Chain.KeplrName != "" {
		sb.WriteString(fmt.Sprintf(
			"<a href='%s'>Keplr</a>\n",
			e.Chain.GetKeplrLink(e.ProposalID),
		))
	}

	explorerLinks := e.Chain.GetExplorerProposalsLinks(e.ProposalID)
	for _, link := range explorerLinks {
		sb.WriteString(fmt.Sprintf(
			"<a href='%s'>%s</a>\n",
			link.Link,
			link.Name,
		))
	}

	return sb.String()
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
	return reporter.BotReply(c, fmt.Sprintf(
		"Notification for proposal #%s on %s are muted till %s.",
		mute.ProposalID,
		mute.Chain,
		mute.Expires.Format(time.RFC1123),
	))
}

func (reporter *TelegramReporter) HandleListMutes(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list mutes query")

	var sb strings.Builder
	sb.WriteString("<strong>Active mutes:</strong>\n\n")

	mutesCount := 0

	for _, mute := range reporter.MutesManager.Mutes.Mutes {
		if mute.IsExpired() {
			continue
		}

		mutesCount++

		sb.WriteString(fmt.Sprintf(
			"<strong>Chain: </strong>%s\n<strong>Proposal ID: </strong>%s\n<strong>Expires: </strong>%s\n\n",
			mute.Chain, mute.ProposalID, mute.Expires,
		))
	}

	if mutesCount == 0 {
		sb.WriteString("No active mutes.")
	}

	return reporter.BotReply(c, sb.String())
}

func (reporter *TelegramReporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	var sb strings.Builder
	sb.WriteString("<strong>cosmos-proposals-checker</strong>\n\n")
	sb.WriteString("Notifies you about the proposals your wallets hasn't voted upon.\n")
	sb.WriteString("Can understand the following commands:\n")
	sb.WriteString("- /proposals_mute &lt;duration&gt; &lt;chain&gt; &lt;proposal ID&gt; - mute notifications for a specific proposal\n")
	sb.WriteString("- /proposals_mutes - display the active proposals mutes list\n")
	sb.WriteString("- /help - display this command\n")
	sb.WriteString("Created by <a href=\"https://freak12techno.github.io\">freak12techno</a> with ❤️.\n")
	sb.WriteString("This bot is open-sourced, you can get the source code at https://github.com/freak12techno/cosmos-proposals-checker.\n\n")
	sb.WriteString("If you like what we're doing, consider <a href=\"https://freak12techno.github.io/validators\">staking with us</a>!\n")

	return reporter.BotReply(c, sb.String())
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
	if len(args) <= 3 {
		return Mute{}, "Usage: /proposals_mute <duration> <chain> <proposal>"
	}

	_, args = args[0], args[1:] // removing first argument as it's always /proposals_mute
	durationString, chain, proposalID := args[0], args[1], args[2]

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
