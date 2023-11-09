package state

import (
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Generator struct {
	Logger zerolog.Logger
	Chains types.Chains
	Mutex  sync.Mutex
}

func NewStateGenerator(logger *zerolog.Logger, chains types.Chains) *Generator {
	return &Generator{
		Logger: logger.With().Str("component", "state_generator").Logger(),
		Chains: chains,
	}
}

func (g *Generator) GetState(oldState State) State {
	state := NewState()

	var wg sync.WaitGroup
	wg.Add(len(g.Chains))

	for _, chain := range g.Chains {
		g.Logger.Info().Str("name", chain.Name).Msg("Processing a chain")
		go func(c *types.Chain) {
			g.ProcessChain(c, state, oldState)
			wg.Done()
		}(chain)
	}

	wg.Wait()
	return state
}

func (g *Generator) ProcessChain(
	chain *types.Chain,
	state State,
	oldState State,
) {
	rpc := tendermint.NewRPC(chain, g.Logger)

	proposals, err := rpc.GetAllProposals()
	if err != nil {
		g.Logger.Warn().Err(err).Msg("Error processing proposals")
		g.Mutex.Lock()
		defer g.Mutex.Unlock()

		state.SetChainProposalsError(chain, err)

		stateChain, found := oldState.ChainInfos[chain.Name]
		if found {
			g.Logger.Trace().Str("chain", chain.Name).Msg("Got older state present, saving it")
			state.SetChainVotes(chain, stateChain.ProposalVotes)
		}

		return
	}

	g.Logger.Info().Int("len", len(proposals)).Msg("Got proposals")

	var wg sync.WaitGroup

	for _, proposal := range proposals {
		g.Logger.Trace().
			Str("name", chain.Name).
			Str("proposal", proposal.ID).
			Msg("Processing a proposal")

		for _, wallet := range chain.Wallets {
			g.Logger.Trace().
				Str("name", chain.Name).
				Str("proposal", proposal.ID).
				Str("wallet", wallet.Address).
				Msg("Processing wallet vote")
			wg.Add(1)

			go func(p types.Proposal, w *types.Wallet) {
				g.ProcessProposalAndWallet(chain, p, rpc, w, state, oldState)
				wg.Done()
			}(proposal, wallet)
		}
	}

	wg.Wait()
}

func (g *Generator) ProcessProposalAndWallet(
	chain *types.Chain,
	proposal types.Proposal,
	rpc *tendermint.RPC,
	wallet *types.Wallet,
	state State,
	oldState State,
) {
	oldVote, _, found := oldState.GetVoteAndProposal(chain.Name, proposal.ID, wallet.Address)
	voteResponse, err := rpc.GetVote(proposal.ID, wallet.Address)

	if found && oldVote.HasVoted() && voteResponse.Vote == nil {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
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

	proposalVote := ProposalVote{
		Wallet: wallet,
	}

	if err != nil {
		proposalVote.Error = types.NewJSONError(err.Error())
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
