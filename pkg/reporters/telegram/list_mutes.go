package telegram

import (
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

	template, err := reporter.TemplatesManager.Render("mutes", filteredMutes)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering template")
		return reporter.BotReply(c, "Error rendering template")
	}

	return reporter.BotReply(c, template)
}
