package report

import (
	"main/pkg/events"
	"main/pkg/report/entry"
	"main/pkg/reporters"
	"main/pkg/state"
	"main/pkg/tendermint"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Generator struct {
	StateManager *state.Manager
	Chains       types.Chains
	RPC          *tendermint.RPC
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

				oldVote, _, _ := oldState.GetVoteAndProposal(chainName, proposalID, wallet)
				newVote, proposal, _ := newState.GetVoteAndProposal(chainName, proposalID, wallet)

				// Error querying for vote - need to notify via Telegram.
				if newVote.IsError() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposalID).
						Str("wallet", wallet).
						Msg("Error querying for vote - sending an alert")
					entries = append(entries, events.VoteQueryError{
						Chain:    g.Chains.FindByName(chainName),
						Proposal: proposal,
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
						Proposal: proposal,
					})
					continue
				}

				// Hasn't voted before but voted now - need to close alert/notify about new vote.
				if newVote.HasVoted() && !oldVote.HasVoted() {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposal.ID).
						Str("wallet", wallet).
						Msg("Wallet hasn't voted before but voted now - closing an alert")

					entries = append(entries, events.VotedEvent{
						Chain:    g.Chains.FindByName(chainName),
						Wallet:   newVote.Wallet,
						Proposal: proposal,
						Vote:     newVote.Vote,
					})
					continue
				}

				// Changed its vote - only notify via Telegram.
				if newVote.HasVoted() && oldVote.HasVoted() && newVote.Vote.Option != oldVote.Vote.Option {
					g.Logger.Debug().
						Str("chain", chainName).
						Str("proposal", proposal.ID).
						Str("wallet", wallet).
						Msg("Wallet changed its vote - sending an alert")

					entries = append(entries, events.RevotedEvent{
						Chain:    g.Chains.FindByName(chainName),
						Wallet:   newVote.Wallet,
						Proposal: proposal,
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
