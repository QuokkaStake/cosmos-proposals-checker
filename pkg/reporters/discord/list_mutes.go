package discord

import (
	mutes "main/pkg/mutes"
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
			filteredMutes := utils.Filter(reporter.MutesManager.Mutes.Mutes, func(m *mutes.Mute) bool {
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
