package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	return reporter.ReplyRender(c, "help", reporter.Version)
}
