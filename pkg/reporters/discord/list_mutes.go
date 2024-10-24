package discord

import (
	"main/pkg/types"
	"main/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetMutesCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "proposals_mutes",
			Description: "List all active mutes.",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mutes, err := reporter.MutesManager.GetAllMutes()
			if err != nil {
				reporter.Logger.Error().Err(err).Msg("Error getting all mutes")
				return
			}

			filteredMutes := utils.Filter(mutes, func(m *types.Mute) bool {
				return !m.IsExpired()
			})

			template, err := reporter.TemplatesManager.Render("mutes", filteredMutes)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "mutes").Msg("Error rendering template")
				return
			}

			reporter.BotRespond(s, i, template)
		},
	}
}
