package telegram

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleTally(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got tally list query")

	if err := reporter.BotReply(c, "Calculating tally for proposals. This might take a while..."); err != nil {
		return err
	}

	tallies, err := reporter.DataManager.GetTallies(context.Background())
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting tallies info: %s", err))
	}

	template, err := reporter.TemplatesManager.Render("tally", tallies)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering template")
		return reporter.BotReply(c, "Error rendering template")
	}

	return reporter.BotReply(c, template)
}
