package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetParamsCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "params",
			Description: "List all chains params.",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			params, err := reporter.DataManager.GetParams(context.Background())
			if err != nil {
				reporter.BotRespond(s, i, fmt.Sprintf("Error getting chain params: %s", err))
				return
			}

			template, err := reporter.TemplatesManager.Render("params", params)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "params").Msg("Error rendering template")
				return
			}

			reporter.BotRespond(s, i, template)
		},
	}
}
