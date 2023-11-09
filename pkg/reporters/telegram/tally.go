package telegram

import (
	"bytes"
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

	tallies, err := reporter.DataManager.GetTallies()
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting tallies info: %s", err))
	}

	template, err := reporter.GetTemplate("tally")
	if err != nil {
		reporter.Logger.Error().
			Err(err).
			Msg("Error rendering tallies template")
		return reporter.BotReply(c, "Error rendering tallies template")
	}

	var buffer bytes.Buffer
	if err := template.Execute(&buffer, tallies); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering votes template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
