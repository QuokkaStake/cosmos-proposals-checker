package telegram

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleParams(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got params query")

	params, err := reporter.DataManager.GetParams(context.Background())
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting chain params: %s", err))
	}

	return reporter.ReplyRender(c, "params", params)
}
