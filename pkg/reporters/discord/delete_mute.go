package discord

import (
	mutes "main/pkg/mutes"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetDeleteMuteCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "proposals_unmute",
			Description: "Unmute proposals' notifications",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "chain",
					Description: "Chain to mute notifications on",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "proposal",
					Description: "Proposal to mute notifications on",
					Required:    false,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			var chain string
			var proposal string

			for _, opt := range options {
				if opt.Name == "chain" {
					chain, _ = opt.Value.(string)
				}
				if opt.Name == "proposal" {
					proposal, _ = opt.Value.(string)
				}
			}

			mute := &mutes.Mute{
				Chain:      chain,
				ProposalID: proposal,
				Expires:    time.Now(),
				Comment:    "",
			}

			if found := reporter.MutesManager.DeleteMute(mute); !found {
				reporter.BotRespond(s, i, "Could not find the mute to delete!")
				return
			}

			template, err := reporter.TemplatesManager.Render("mute_deleted", mute)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "mute_deleted").Msg("Error rendering template")
				return
			}

			reporter.BotRespond(s, i, template)
		},
	}
}
