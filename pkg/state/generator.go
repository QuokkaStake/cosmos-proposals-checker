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

	g.Mutex.Lock()
	state.SetChainProposalsHeight(chain, proposalsHeight)
	g.Mutex.Unlock()

	var wg sync.WaitGroup

	for _, proposal := range proposals {
		wg.Add(1)

		go func(p types.Proposal) {
			g.ProcessProposal(chain, p, fetcher, state, oldState)
			wg.Done()
		}(proposal)
	}

	wg.Wait()
}

func (g *Generator) ProcessProposal(
	chain *types.Chain,
	proposal types.Proposal,
	fetcher fetchersPkg.Fetcher,
	state State,
	oldState State,
) {
	g.Logger.Trace().
		Str("name", chain.Name).
		Str("proposal", proposal.ID).
		Msg("Processing a proposal")

	if !proposal.IsInVoting() {
		g.Logger.Trace().
			Str("name", chain.Name).
			Str("proposal", proposal.ID).
			Msg("Proposal is not in voting period - not processing it")
		g.Mutex.Lock()
		state.SetProposal(chain, proposal)
		g.Mutex.Unlock()
		return
	}

	var wg sync.WaitGroup

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
	oldVote, found := oldState.GetVote(chain.Name, proposal.ID, wallet.Address)
	vote, voteHeight, err := fetcher.GetVote(proposal.ID, wallet.Address, oldVote.Height)

	proposalVote := ProposalVote{
		Wallet: wallet,
	}

	if err != nil {
		// 1. If error occurred - store the error, but preserve the older height and vote.
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Str("wallet", wallet.Address).
			Int64("height", voteHeight).
			Err(err).
			Msg("Error fetching wallet vote - preserving the older height and vote")

		proposalVote.Error = err
		if found {
			proposalVote.Height = oldVote.Height
			proposalVote.Vote = oldVote.Vote
		}
	} else if found && oldVote.HasVoted() && vote == nil {
		// 2. If there's no newer vote while there's an older vote - preserve the older vote
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Str("wallet", wallet.Address).
			Msg("Wallet has voted and there's no vote in the new state - using old vote")

		proposalVote.Vote = oldVote.Vote
		proposalVote.Height = voteHeight
	} else {
		// 3. Wallet voted (or hadn't voted and hadn't voted before) - use the older vote.
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
