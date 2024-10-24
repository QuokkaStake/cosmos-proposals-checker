package discord

import (
	"fmt"
	"main/pkg/types"
	"time"

	"github.com/guregu/null/v5"

	"github.com/bwmarrin/discordgo"
)

func (reporter *Reporter) GetAddMuteCommand() *Command {
	return &Command{
		Info: &discordgo.ApplicationCommand{
			Name:        "proposals_mute",
			Description: "Mute proposals' notifications",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "duration",
					Description: "For how long to mute notifications",
					Required:    true,
				},
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

			durationString, _ := options[0].Value.(string)
			chain := null.NewString("", false)
			proposal := null.NewString("", false)

			_, opts := options[0], options[1:]

			for _, opt := range opts {
				if opt.Name == "chain" {
					chainRaw, _ := opt.Value.(string)
					chain = null.StringFrom(chainRaw)
				}
				if opt.Name == "proposal" {
					proposalRaw, _ := opt.Value.(string)
					proposal = null.StringFrom(proposalRaw)
				}
			}

			duration, err := time.ParseDuration(durationString)
			if err != nil {
				reporter.BotRespond(s, i, "Invalid mute duration provided: %s")
				return
			}

			mute := &types.Mute{
				Chain:      chain,
				ProposalID: proposal,
				Expires:    time.Now().Add(duration),
				Comment: fmt.Sprintf(
					"Muted using cosmos-proposals-checker for %s",
					duration,
				),
			}

			if insertErr := reporter.MutesManager.AddMute(mute); insertErr != nil {
				reporter.Logger.Error().Err(err).Msg("Error adding mute")
				reporter.BotRespond(s, i, "Error adding mute!")
				return
			}

			template, err := reporter.TemplatesManager.Render("mute_added", mute)
			if err != nil {
				reporter.Logger.Error().Err(err).Str("template", "mute_added").Msg("Error rendering template")
				return
			}

			reporter.BotRespond(s, i, template)
		},
	}
}
