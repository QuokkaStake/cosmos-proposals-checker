package telegram

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleParams(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got params query")

	params, err := reporter.DataManager.GetParams()
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting chain params: %s", err))
	}

	template, err := reporter.TemplatesManager.Render("params", params)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering template")
		return reporter.BotReply(c, "Error rendering template")
	}

	return reporter.BotReply(c, template)
}
