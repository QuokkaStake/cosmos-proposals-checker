package telegram

import (
	"bytes"
	statePkg "main/pkg/state"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleProposals(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got proposals list query")

	state := reporter.StateGenerator.GetState(statePkg.NewState())
	renderedState := state.ToRenderedState()

	template, _ := reporter.GetTemplate("proposals")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, renderedState); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering proposals template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
