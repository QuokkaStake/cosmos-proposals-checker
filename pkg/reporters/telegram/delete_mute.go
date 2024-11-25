package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleDeleteMute(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got delete mute query")

	mute, err := ParseMuteDeleteOptions(c)
	if err != "" {
		return c.Reply("Error deleting mute: " + err)
	}

	if found, deleteErr := reporter.MutesManager.DeleteMute(mute); !found {
		return c.Reply("Could not find the mute to delete!")
	} else if deleteErr != nil {
		return c.Reply("Error deleting mute!")
	}

	return reporter.ReplyRender(c, "mute_deleted", mute)
}
