package discord

import (
	"context"
	"fmt"
	"main/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetTallyCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "tally",
			Description: "Get active proposals' tallies",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			reporter.BotSendInteraction(s, i, "Calculating tally for proposals. This might take a while...")

			tallies, err := reporter.DataManager.GetTallies(context.Background())
			if err != nil {
				reporter.BotRespond(s, i, fmt.Sprintf("Error getting tallies info: %s", err))
				return
			}

			template, err := reporter.TemplatesManager.Render("tally", tallies)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "tally").Msg("Error rendering template")
				return
			}

			chunks := utils.SplitStringIntoChunks(template, 2000)
			for _, chunk := range chunks {
				reporter.BotSendFollowup(s, i, chunk)
			}
		},
	}
}
