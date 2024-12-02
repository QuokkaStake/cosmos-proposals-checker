package telegram

import (
	"fmt"

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

	if found, deleteErr := reporter.MutesManager.DeleteMute(mute); deleteErr != nil {
		return c.Reply(fmt.Sprintf("Error deleting mute: %s!", deleteErr))
	} else if !found {
		return c.Reply("Could not find the mute to delete!")
	}

	return reporter.ReplyRender(c, "mute_deleted", mute)
}
