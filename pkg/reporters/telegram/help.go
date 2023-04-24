package telegram

import (
	"bytes"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	template, _ := reporter.GetTemplate("help")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, nil); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering help template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
