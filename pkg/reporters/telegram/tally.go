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

	msg, err := reporter.TelegramBot.Reply(c.Message(), "Calculating tally for proposals. This might take a while...")
	if err != nil {
		return err
	}

	tallies, err := reporter.DataManager.GetTallies(context.Background())
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting tallies info: %s", err))
	}

	return reporter.EditRender(c, msg, "tally", tallies)
}
