package main

import (
	"time"

	"github.com/rs/zerolog"
)

type ReportGenerator struct {
	StateManager *StateManager
	Chains       Chains
	RPC          *RPC
	Logger       zerolog.Logger
}

func NewReportGenerator(
	manager *StateManager,
	logger *zerolog.Logger,
	chains Chains,
) *ReportGenerator {
	return &ReportGenerator{
		StateManager: manager,
		Chains:       chains,
		Logger:       logger.With().Str("component", "report_generator").Logger(),
	}
}

func (g *ReportGenerator) GenerateReport(oldState, newState State) Report {
	entries := []ReportEntry{}

	for chainName, chainInfo := range newState.ChainInfos {
		if chainInfo.HasProposalsError() {
			g.Logger.Debug().
				Str("chain", chainName).
				Msg("Error querying for proposals - sending an alert")
			entry := ReportEntry{
				Chain:                  *g.Chains.FindByName(chainName),
				Wallet:                 "",
				ProposalID:             "",
				ProposalTitle:          "",
				ProposalDescription:    "",
				ProposalVoteEndingTime: time.Now(),
				Type:                   ProposalQueryError,
				Value:                  chainInfo.ProposalsError,
			}
			entries = append(entries, entry)
			continue
		}

		for proposalID, proposalVotes := range chainInfo.ProposalVotes {
			for wallet := range proposalVotes.Votes {
				g.Logger.Trace().
					Str("name", chainName).
					Str("proposal", proposalID).
					Str("wallet", wallet).
					Msg("Generating report for a wallet vote")

				oldVote, _, _ := oldState.GetVoteAndProposal(chainName, proposalID, wallet)
				newVote, proposal, _ := newState.GetVoteAndProposal(chainName, proposalID, wallet)

				entry := ReportEntry{
					Chain:                  *g.Chains.FindByName(chainName),
					Wallet:                 wallet,
					ProposalID:             proposalID,
					ProposalTitle:          proposal.Content.Title,
					ProposalDescription:    proposal.Content.Description,
					ProposalVoteEndingTime: proposal.VotingEndTime,
				}

				// Error querying for vote - need to notify via Telegram.
				if newVote.IsError() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposalID).
						Str("wallet", wallet).
						Msg("Error querying for vote - sending an alert")
					entry.Type = VoteQueryError
					entry.Value = newVote.Error

					entries = append(entries, entry)
					continue
				}

				// Hasn't voted for this proposal - need to notify.
				if !newVote.HasVoted() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposalID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted now - sending an alert")
					entry.Type = NotVoted
					entries = append(entries, entry)
					continue
				}

				// Hasn't voted before but voted now - need to close alert/notify about new vote.
				if newVote.HasVoted() && !oldVote.HasVoted() {
					vote := *newVote.Vote

					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted before but voted now - closing an alert")

					entry.Type = Voted
					entry.Value = vote.Option

					entries = append(entries, entry)
					continue
				}

				// Changed its vote - only notify via Telegram.
				if newVote.HasVoted() && oldVote.HasVoted() && newVote.Vote.Option != oldVote.Vote.Option {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet).
						Msg("Wallet changed its vote - sending an alert")
					entry.Type = Revoted
					entry.Value = newVote.Vote.Option
					entry.OldValue = oldVote.Vote.Option

					entries = append(entries, entry)
				}
			}
		}
	}

	g.StateManager.CommitState(newState)

	return Report{Entries: entries}
}
