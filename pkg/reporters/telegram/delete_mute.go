package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleDeleteMute(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got delete mute query")

	mute, err := ParseMuteDeleteOptions(c.Text(), c)
	if err != "" {
		return c.Reply("Error deleting mute: " + err)
	}

	if found, deleteErr := reporter.MutesManager.DeleteMute(mute); !found {
		return c.Reply("Could not find the mute to delete!")
	} else if deleteErr != nil {
		return c.Reply("Error deleting mute!")
	}

	templateRendered, renderErr := reporter.TemplatesManager.Render("mute_deleted", mute)
	if renderErr != nil {
		reporter.Logger.Error().Err(renderErr).Msg("Error rendering template")
		return reporter.BotReply(c, "Error rendering template")
	}

	return reporter.BotReply(c, templateRendered)
}
