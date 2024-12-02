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

	return reporter.ReplyRender(c, "proposals", renderedState)
}
