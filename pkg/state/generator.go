package state

import (
	fetchersPkg "main/pkg/fetchers"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Generator struct {
	Logger   zerolog.Logger
	Chains   types.Chains
	Fetchers map[string]fetchersPkg.Fetcher
	Mutex    sync.Mutex
}

func NewStateGenerator(logger *zerolog.Logger, chains types.Chains) *Generator {
	fetchers := make(map[string]fetchersPkg.Fetcher, len(chains))

	for _, chain := range chains {
		fetchers[chain.Name] = fetchersPkg.GetFetcher(chain, *logger)
	}

	return &Generator{
		Logger:   logger.With().Str("component", "state_generator").Logger(),
		Chains:   chains,
		Fetchers: fetchers,
	}
}

func (g *Generator) GetState(oldState State) State {
	state := NewState()

	var wg sync.WaitGroup
	wg.Add(len(g.Chains))

	for _, chain := range g.Chains {
		g.Logger.Info().Str("name", chain.Name).Msg("Processing a chain")

		fetcher := g.Fetchers[chain.Name]

		go func(c *types.Chain) {
			g.ProcessChain(c, state, oldState, fetcher)
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
	fetcher fetchersPkg.Fetcher,
) {
	prevHeight := oldState.GetLastProposalsHeight(chain)
	proposals, proposalsHeight, err := fetcher.GetAllProposals(prevHeight)
	if err != nil {
		g.Logger.Warn().Err(err).Msg("Error processing proposals")
		g.Mutex.Lock()
		defer g.Mutex.Unlock()

		state.SetChainProposalsError(chain, err)
		state.SetChainProposalsHeight(chain, prevHeight)

		stateChain, found := oldState.ChainInfos[chain.Name]
		if found {
			g.Logger.Trace().Str("chain", chain.Name).Msg("Got older state present, saving it")
			state.SetChainVotes(chain, stateChain.ProposalVotes)
		}

		return
	}

	g.Logger.Info().
		Str("chain", chain.Name).
		Int("len", len(proposals)).
		Int64("height", proposalsHeight).
		Msg("Got proposals")

	state.SetChainProposalsHeight(chain, proposalsHeight)

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
				g.ProcessProposalAndWallet(chain, p, fetcher, w, state, oldState)
				wg.Done()
			}(proposal, wallet)
		}
	}

	wg.Wait()
}

func (g *Generator) ProcessProposalAndWallet(
	chain *types.Chain,
	proposal types.Proposal,
	fetcher fetchersPkg.Fetcher,
	wallet *types.Wallet,
	state State,
	oldState State,
) {
	oldVote, _, found := oldState.GetVoteAndProposal(chain.Name, proposal.ID, wallet.Address)
	vote, voteHeight, err := fetcher.GetVote(proposal.ID, wallet.Address, oldVote.Height)

	if found && oldVote.HasVoted() && vote == nil {
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
		proposalVote.Error = err
		if found {
			proposalVote.Height = oldVote.Height
		}
	} else {
		proposalVote.Vote = vote
		proposalVote.Height = voteHeight
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
