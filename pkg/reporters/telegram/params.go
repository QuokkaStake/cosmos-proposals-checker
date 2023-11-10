package telegram

import (
	"bytes"
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

	template, err := reporter.GetTemplate("params")
	if err != nil {
		reporter.Logger.Error().
			Err(err).
			Msg("Error rendering params template")
		return reporter.BotReply(c, "Error rendering params template")
	}

	var buffer bytes.Buffer
	if err := template.Execute(&buffer, params); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering params template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
