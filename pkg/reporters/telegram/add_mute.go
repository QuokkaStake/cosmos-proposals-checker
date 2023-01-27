package telegram

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func (reporter *TelegramReporter) HandleAddMute(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got add mute query")

	mute, err := ParseMuteOptions(c.Text(), c)
	if err != "" {
		return c.Reply(fmt.Sprintf("Error muting notification: %s", err))
	}

	reporter.MutesManager.AddMute(mute)
	if mute.ProposalID != "" {
		return reporter.BotReply(c, fmt.Sprintf(
			"Notification for proposal #%s on %s are muted till %s.",
			mute.ProposalID,
			mute.Chain,
			mute.GetExpirationTime(),
		))
	}

	return reporter.BotReply(c, fmt.Sprintf(
		"Notification for all proposals on %s are muted till %s.",
		mute.Chain,
		mute.GetExpirationTime(),
	))
}
