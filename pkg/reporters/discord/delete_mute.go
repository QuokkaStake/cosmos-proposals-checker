package discord

import (
	"main/pkg/types"
	"time"

	"github.com/guregu/null/v5"

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

			chain := null.NewString("", false)
			proposal := null.NewString("", false)

			for _, opt := range options {
				if opt.Name == "chain" {
					chainRaw, _ := opt.Value.(string)
					chain = null.StringFrom(chainRaw)
				}
				if opt.Name == "proposal" {
					proposalRaw, _ := opt.Value.(string)
					proposal = null.StringFrom(proposalRaw)
				}
			}

			mute := &types.Mute{
				Chain:      chain,
				ProposalID: proposal,
				Expires:    time.Now(),
				Comment:    "",
			}

			if found, insertErr := reporter.MutesManager.DeleteMute(mute); !found {
				reporter.BotRespond(s, i, "Could not find the mute to delete!")
				return
			} else if insertErr != nil {
				reporter.BotRespond(s, i, "Error deleting mute!")
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
