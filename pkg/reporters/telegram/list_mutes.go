package telegram

import (
	"bytes"
	mutes "main/pkg/mutes"
	"main/pkg/utils"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleListMutes(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list mutes query")

	filteredMutes := utils.Filter(reporter.MutesManager.Mutes.Mutes, func(m *mutes.Mute) bool {
		return !m.IsExpired()
	})

	template, _ := reporter.GetTemplate("mutes")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, filteredMutes); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering votes template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
