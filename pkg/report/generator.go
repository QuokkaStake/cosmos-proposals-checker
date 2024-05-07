package report

import (
	"main/pkg/events"
	"main/pkg/fetchers/cosmos"
	"main/pkg/report/entry"
	"main/pkg/reporters"
	"main/pkg/state"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Generator struct {
	StateManager *state.Manager
	Chains       types.Chains
	RPC          *cosmos.RPC
	Logger       zerolog.Logger
}

func NewReportGenerator(
	manager *state.Manager,
	logger *zerolog.Logger,
	chains types.Chains,
) *Generator {
	return &Generator{
		StateManager: manager,
		Chains:       chains,
		Logger:       logger.With().Str("component", "report_generator").Logger(),
	}
}

func (g *Generator) GenerateReport(oldState, newState state.State) reporters.Report {
	entries := []entry.ReportEntry{}

	for chainName, chainInfo := range newState.ChainInfos {
		if chainInfo.HasProposalsError() {
			g.Logger.Debug().
				Str("chain", chainName).
				Msg("Error querying for proposals - sending an alert")
			entries = append(entries, events.ProposalsQueryErrorEvent{
				Chain: g.Chains.FindByName(chainName),
				Error: chainInfo.ProposalsError,
			})
			continue
		}

		for proposalID, proposalVotes := range chainInfo.ProposalVotes {
			for wallet := range proposalVotes.Votes {
				g.Logger.Trace().
					Str("name", chainName).
					Str("proposal", proposalID).
					Str("wallet", wallet).
					Msg("Generating report for a wallet vote")

				oldVote, _ := oldState.GetVote(chainName, proposalID, wallet)
				newVote, _ := newState.GetVote(chainName, proposalID, wallet)

				oldProposal, oldProposalFound := oldState.GetProposal(chainName, proposalID)
				newProposal := proposalVotes.Proposal

				if oldProposalFound {
					// There can be the following cases:
					// 1) old proposal not found - this is a new proposal, no custom logic needed.
					// 2) old proposal is found, it's not in voting and the current one is not in voting
					// (like deposit -> deposit) - we don't need to report it at all.
					// 3) old proposal is found, it's not in voting but the current one is in voting
					// (like deposit -> deposit) - no custom logic needed.
					// 4) old proposal is found, it's in voting and the current one is in voting
					// (like voting -> voting) - no custom logic needed.
					// 5) old proposal is found, it's in voting and the current one is not in voting
					// (like voting -> passed) -> send a new event that the voting has finished.

					// case 1, 3 and 4 is handled outside of this if case
					if oldProposal.IsInVoting() && !newProposal.IsInVoting() { // case 5
						g.Logger.Debug().
							Str("chain", chainName).
							Str("proposal", proposalID).
							Str("wallet", wallet).
							Msg("Previously proposal was in voting, but it's not now - sending an alert")

						entries = append(entries, events.FinishedVotingEvent{
							Chain:    g.Chains.FindByName(chainName),
							Proposal: newProposal,
						})
						continue
					} else if !oldProposal.IsInVoting() && !newProposal.IsInVoting() { // case 2
						g.Logger.Debug().
							Str("chain", chainName).
							Str("proposal", proposalID).
							Str("wallet", wallet).
							Msg("Previously proposal was and is not in voting period - ignoring.")
						continue
					}
				}

				// Error querying for vote - need to notify via Telegram.
				if newVote.IsError() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposalID).
						Str("wallet", wallet).
						Msg("Error querying for vote - sending an alert")
					entries = append(entries, events.VoteQueryError{
						Chain:    g.Chains.FindByName(chainName),
						Proposal: newProposal,
						Error:    newVote.Error,
					})

					continue
				}

				// Hasn't voted for this proposal - need to notify.
				if !newVote.HasVoted() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposalID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted now - sending an alert")
					entries = append(entries, events.NotVotedEvent{
						Chain:    g.Chains.FindByName(chainName),
						Wallet:   newVote.Wallet,
						Proposal: newProposal,
					})
					continue
				}

				// Hasn't voted before but voted now - need to close alert/notify about new vote.
				if newVote.HasVoted() && !oldVote.HasVoted() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", newProposal.ID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted before but voted now - closing an alert")

					entries = append(entries, events.VotedEvent{
						Chain:    g.Chains.FindByName(chainName),
						Wallet:   newVote.Wallet,
						Proposal: newProposal,
						Vote:     newVote.Vote,
					})
					continue
				}

				// Changed its vote - only notify via Telegram.
				if newVote.HasVoted() && oldVote.HasVoted() && !newVote.Vote.VotesEquals(oldVote.Vote) {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", newProposal.ID).
						Str("wallet", wallet).
						Msg("Wallet changed its vote - sending an alert")

					entries = append(entries, events.RevotedEvent{
						Chain:    g.Chains.FindByName(chainName),
						Wallet:   newVote.Wallet,
						Proposal: newProposal,
						Vote:     newVote.Vote,
						OldVote:  oldVote.Vote,
					})
				}
			}
		}
	}

	g.StateManager.CommitState(newState)

	return reporters.Report{Entries: entries}
}
