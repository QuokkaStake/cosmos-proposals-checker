package telegram

import (
	"context"
	statePkg "main/pkg/state"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleProposals(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got proposals list query")

	state := reporter.StateGenerator.GetState(statePkg.NewState(), context.Background())
	renderedState := state.ToRenderedState()

	template, err := reporter.TemplatesManager.Render("proposals", renderedState)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering template")
		return reporter.BotReply(c, "Error rendering template")
	}

	return reporter.BotReply(c, template)
}
