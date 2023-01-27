package telegram

import (
	"bytes"
	tele "gopkg.in/telebot.v3"
	"main/pkg/state"
)

func (reporter *TelegramReporter) HandleProposals(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got proposals list query")

	state := reporter.StateGenerator.GetState(state.NewState())
	template, _ := reporter.GetTemplate("proposals")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, state); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering votes template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
