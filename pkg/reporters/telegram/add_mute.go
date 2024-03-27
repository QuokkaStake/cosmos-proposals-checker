package telegram

import (
	"bytes"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleAddMute(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got add mute query")

	mute, err := ParseMuteOptions(c.Text(), c)
	if err != "" {
		return c.Reply("Error muting notification: " + err)
	}

	reporter.MutesManager.AddMute(mute)

	template, _ := reporter.GetTemplate("mute_added")
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, mute); err != nil {
		reporter.Logger.Error().Err(err).Msg("Error rendering mute_added template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
