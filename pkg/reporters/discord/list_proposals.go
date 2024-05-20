package discord

import (
	"context"
	statePkg "main/pkg/state"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetProposalsCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "proposals",
			Description: "Get list of active proposals and your wallet's votes on them.",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			state := reporter.StateGenerator.GetState(statePkg.NewState(), context.Background())
			renderedState := state.ToRenderedState()

			template, err := reporter.TemplatesManager.Render("proposals", renderedState)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "proposals").Msg("Error rendering template")
				return
			}

			reporter.BotRespond(s, i, template)
		},
	}
}
