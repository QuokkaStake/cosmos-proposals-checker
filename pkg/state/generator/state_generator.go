package generator

import (
	configTypes "main/pkg/config/types"
	statePackage "main/pkg/state"
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type StateGenerator struct {
	Logger zerolog.Logger
	Chains configTypes.Chains
	Mutex  sync.Mutex
}

func NewStateGenerator(logger *zerolog.Logger, chains configTypes.Chains) *StateGenerator {
	return &StateGenerator{
		Logger: logger.With().Str("component", "state_generator").Logger(),
		Chains: chains,
	}
}

func (g *StateGenerator) GetState(oldState statePackage.State) statePackage.State {
	state := statePackage.NewState()

	var wg sync.WaitGroup
	wg.Add(len(g.Chains))

	for _, chain := range g.Chains {
		g.Logger.Info().Str("name", chain.Name).Msg("Processing a chain")
		go func(c *configTypes.Chain) {
			g.ProcessChain(c, state, oldState)
			wg.Done()
		}(chain)
	}

	wg.Wait()
	return state
}

func (g *StateGenerator) ProcessChain(
	chain *configTypes.Chain,
	state statePackage.State,
	oldState statePackage.State,
) {
	rpc := tendermint.NewRPC(chain.LCDEndpoints, g.Logger)

	proposals, err := rpc.GetAllProposals()
	if err != nil {
		g.Logger.Warn().Err(err).Msg("Error processing proposals")
		g.Mutex.Lock()
		state.SetChainProposalsError(chain, err)
		g.Mutex.Unlock()
		return
	}

	g.Logger.Info().Int("len", len(proposals)).Msg("Got proposals")

	var wg sync.WaitGroup

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
			wg.Add(1)

			go func(p types.Proposal, w *configTypes.Wallet) {
				g.ProcessProposalAndWallet(chain, p, rpc, w, state, oldState)
				wg.Done()
			}(proposal, wallet)
		}
	}

	wg.Wait()
}

func (g *StateGenerator) ProcessProposalAndWallet(
	chain *configTypes.Chain,
	proposal types.Proposal,
	rpc *tendermint.RPC,
	wallet *configTypes.Wallet,
	state statePackage.State,
	oldState statePackage.State,
) {
	oldVote, _, found := oldState.GetVoteAndProposal(chain.Name, proposal.ProposalID, wallet.Address)
	voteResponse, err := rpc.GetVote(proposal.ProposalID, wallet.Address)

	if found && oldVote.HasVoted() && voteResponse.Vote == nil {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ProposalID).
			Str("wallet", wallet.Address).
			Msg("Wallet has voted and there's no vote in the new state - using old vote")

		g.Mutex.Lock()
		state.SetVote(
			chain,
			proposal,
			wallet,
			oldVote,
		)
		g.Mutex.Unlock()
	}

	proposalVote := statePackage.ProposalVote{
		Wallet: wallet,
	}

	if err != nil {
		proposalVote.Error = err.Error()
	} else {
		proposalVote.Vote = voteResponse.Vote
	}

	g.Mutex.Lock()
	state.SetVote(
		chain,
		proposal,
		wallet,
		proposalVote,
	)
	g.Mutex.Unlock()
}
