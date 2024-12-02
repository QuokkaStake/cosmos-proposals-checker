package telegram

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleListMutes(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list mutes query")

	mutes, err := reporter.MutesManager.GetAllMutes()
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Error fetching mutes")
		return reporter.BotReply(c, fmt.Sprintf("Error fetching mutes: %s", err))
	}

	filteredMutes := utils.Filter(mutes, func(m *types.Mute) bool {
		return !m.IsExpired()
	})

	return reporter.ReplyRender(c, "mutes", filteredMutes)
}
