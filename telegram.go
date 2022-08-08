package main

import (
	"fmt"
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
}

const MaxMessageSize = 4096

func NewTelegramReporter(config TelegramConfig, mutesManager *MutesManager, logger *zerolog.Logger) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken: config.TelegramToken,
		TelegramChat:  config.TelegramChat,
		MutesManager:  mutesManager,
		Logger:        logger.With().Str("component", "telegram_reporter").Logger(),
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

func (reporter *TelegramReporter) SerializeReportEntry(e ReportEntry) string {
	var sb strings.Builder

	messageText := "🔴 <strong>Wallet %s hasn't voted on proposal %s on %s</strong>\n%s\n"
	if e.HasVoted() {
		messageText = "✅ <strong>Wallet %s has voted on proposal %s on %s</strong>\n%s\n"
	}

	sb.WriteString(fmt.Sprintf(
		messageText,
		e.Wallet,
		e.ProposalID,
		e.Chain.GetName(),
		e.ProposalTitle,
	))

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

	sb.WriteString(
		"\nSent by <a href='https://github.com/freak12techno/cosmos-proposals-checker'>cosmos-proposals-checker.</a>",
	)

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

		serializedEntry := reporter.SerializeReportEntry(entry)

		_, err := reporter.TelegramBot.Send(
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
