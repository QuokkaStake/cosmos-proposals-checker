package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	telegramBot "gopkg.in/tucnak/telebot.v2"
)

type TelegramReporter struct {
	TelegramToken string
	TelegramChat  int64

	TelegramBot *telegramBot.Bot
	Logger      zerolog.Logger
}

func NewTelegramReporter(config TelegramConfig, logger *zerolog.Logger) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken: config.TelegramToken,
		TelegramChat:  config.TelegramChat,
		Logger:        logger.With().Str("component", "telegram_reporter").Logger(),
	}
}

func (reporter *TelegramReporter) Init() {
	if reporter.TelegramToken == "" || reporter.TelegramChat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return
	}

	bot, err := telegramBot.NewBot(telegramBot.Settings{
		Token:  reporter.TelegramToken,
		Poller: &telegramBot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return
	}

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
		serializedEntry := reporter.SerializeReportEntry(entry)

		_, err := reporter.TelegramBot.Send(
			&telegramBot.User{
				ID: reporter.TelegramChat,
			},
			serializedEntry,
			telegramBot.ModeHTML,
			telegramBot.NoPreview,
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
