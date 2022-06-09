package main

import (
	"fmt"

	"github.com/rs/zerolog"
)

type ReportGenerator struct {
	StateManager *StateManager
	Chains       []Chain
	RPC          *RPC
	Logger       zerolog.Logger
}

func NewReportGenerator(
	manager *StateManager,
	logger *zerolog.Logger,
	chains []Chain,
) *ReportGenerator {
	return &ReportGenerator{
		StateManager: manager,
		Chains:       chains,
		Logger:       logger.With().Str("component", "report_generator").Logger(),
	}
}

func (g *ReportGenerator) GenerateReport() Report {
	votesMap := make(map[string]map[string]map[string]*Vote)
	proposalsMap := make(map[string][]Proposal)

	for _, chain := range g.Chains {
		votesMap[chain.Name] = make(map[string]map[string]*Vote)

		rpc := NewRPC(chain.LCDEndpoints, g.Logger)

		g.Logger.Info().Str("name", chain.Name).Msg("Processing a chain")
		proposals, err := rpc.GetAllProposals()
		if err != nil {
			g.Logger.Warn().Err(err).Msg("Error processing proposals")
			continue
		}

		g.Logger.Info().Int("len", len(proposals)).Msg("Got proposals")
		proposalsMap[chain.Name] = proposals

		for _, proposal := range proposals {
			for _, wallet := range chain.Wallets {
				if g.StateManager.HasVotedBefore(chain.Name, proposal.ProposalID, wallet) {
					g.Logger.Trace().
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet).
						Msg("Wallet has already voted, not checking again,")
					g.StateManager.SetVote(
						chain.Name,
						proposal.ProposalID,
						wallet,
						g.StateManager.GetVoteBefore(chain.Name, proposal.ProposalID, wallet),
					)
					continue
				}

				g.Logger.Info().
					Str("proposal", proposal.ProposalID).
					Str("wallet", wallet).
					Msg("Checking if a wallet had voted")

				vote, err := rpc.GetVote(proposal.ProposalID, wallet)
				if err != nil {
					g.Logger.Warn().Err(err).Msg("Error processing vote")
				}

				g.Logger.Info().Str("result", fmt.Sprintf("%+v", vote)).Msg("Got vote")
				g.StateManager.SetVote(chain.Name, proposal.ProposalID, wallet, vote.Vote)
			}
		}
	}

	entries := []ReportEntry{}

	for _, chain := range g.Chains {
		for _, proposal := range proposalsMap[chain.Name] {
			for _, wallet := range chain.Wallets {
				votedNow := g.StateManager.HasVotedNow(chain.Name, proposal.ProposalID, wallet)
				votedBefore := g.StateManager.HasVotedBefore(chain.Name, proposal.ProposalID, wallet)

				// Hasn't voted for this proposal - need to notify.
				if !votedNow {
					g.Logger.Debug().
						Str("chain", chain.Name).
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted now - sending an alert")
					entries = append(entries, ReportEntry{
						Chain:               chain,
						Wallet:              wallet,
						ProposalID:          proposal.ProposalID,
						ProposalTitle:       proposal.Content.Title,
						ProposalDescription: proposal.Content.Description,
					})
				}

				// Hasn't voted before but voted now - need to close alert/notify about new vote.
				if votedNow && !votedBefore {
					vote := g.StateManager.GetVote(chain.Name, proposal.ProposalID, wallet)
					if vote == nil {
						g.Logger.Info().
							Str("chain", chain.Name).
							Str("proposal", proposal.ProposalID).
							Str("wallet", wallet).
							Msg("No vote found while there should be one")
						continue
					}

					g.Logger.Debug().
						Str("chain", chain.Name).
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted before but voted now - closing an alert")

					entries = append(entries, ReportEntry{
						Chain:               chain,
						Wallet:              wallet,
						ProposalID:          proposal.ProposalID,
						ProposalTitle:       proposal.Content.Title,
						ProposalDescription: proposal.Content.Description,
						Vote:                vote.Option,
					})
				}
			}
		}
	}

	g.StateManager.CommitNewState()

	return Report{Entries: entries}
}
