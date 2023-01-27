package generator

import (
	configTypes "main/pkg/config/types"
	statePackage "main/pkg/state"
	"main/pkg/tendermint"

	"github.com/rs/zerolog"
)

type StateGenerator struct {
	Logger zerolog.Logger
	Chains configTypes.Chains
}

func NewStateGenerator(logger *zerolog.Logger, chains configTypes.Chains) *StateGenerator {
	return &StateGenerator{
		Logger: logger.With().Str("component", "state_generator").Logger(),
		Chains: chains,
	}
}

func (g *StateGenerator) GetState(oldState statePackage.State) statePackage.State {
	state := statePackage.NewState()

	for _, chain := range g.Chains {
		g.Logger.Info().Str("name", chain.Name).Msg("Processing a chain")

		rpc := tendermint.NewRPC(chain.LCDEndpoints, g.Logger)

		proposals, err := rpc.GetAllProposals()
		if err != nil {
			g.Logger.Warn().Err(err).Msg("Error processing proposals")
			state.SetChainProposalsError(chain, err)
			continue
		}

		g.Logger.Info().Int("len", len(proposals)).Msg("Got proposals")

		for _, proposal := range proposals {
			g.Logger.Trace().
				Str("name", chain.Name).
				Str("proposal", proposal.ProposalID).
				Msg("Processing a proposal")

			for _, wallet := range chain.Wallets {
				g.Logger.Trace().
					Str("name", chain.Name).
					Str("proposal", proposal.ProposalID).
					Str("wallet", wallet.Address).
					Msg("Processing wallet vote")

				oldVote, _, found := oldState.GetVoteAndProposal(chain.Name, proposal.ProposalID, wallet.Address)
				voteResponse, err := rpc.GetVote(proposal.ProposalID, wallet.Address)

				if found && oldVote.HasVoted() && voteResponse.Vote == nil {
					g.Logger.Trace().
						Str("chain", chain.Name).
						Str("proposal", proposal.ProposalID).
						Str("wallet", wallet.Address).
						Msg("Wallet has voted and there's no vote in the new state - using old vote")

					state.SetVote(
						chain,
						proposal,
						wallet,
						oldVote,
					)

					continue
				}

				proposalVote := statePackage.ProposalVote{
					Wallet: wallet,
				}

				if err != nil {
					proposalVote.Error = err.Error()
				} else {
					proposalVote.Vote = voteResponse.Vote
				}

				state.SetVote(
					chain,
					proposal,
					wallet,
					proposalVote,
				)
			}
		}
	}

	return state
}
