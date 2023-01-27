package telegram

import (
	"bytes"
	tele "gopkg.in/telebot.v3"
	"main/pkg/types"
	"main/pkg/utils"
)

func (reporter *TelegramReporter) HandleListMutes(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list mutes query")

	mutes := utils.Filter(reporter.MutesManager.Mutes.Mutes, func(m types.Mute) bool {
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
