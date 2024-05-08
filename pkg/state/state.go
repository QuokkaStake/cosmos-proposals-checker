package state

import (
	"main/pkg/types"
	"sort"
	"strconv"
)

type ProposalVote struct {
	Wallet *types.Wallet
	Vote   *types.Vote
	Error  *types.QueryError
	Height int64
}

func (v ProposalVote) HasVoted() bool {
	return v.Vote != nil && v.Error == nil
}

func (v ProposalVote) IsError() bool {
	return v.Error != nil
}

type WalletVotes struct {
	Proposal types.Proposal
	Votes    map[string]ProposalVote
}

type ChainInfo struct {
	Chain           *types.Chain
	ProposalVotes   map[string]WalletVotes
	ProposalsError  *types.QueryError
	ProposalsHeight int64
}

func (c ChainInfo) HasProposalsError() bool {
	return c.ProposalsError != nil
}

type State struct {
	ChainInfos map[string]*ChainInfo
}

func NewState() State {
	return State{
		ChainInfos: make(map[string]*ChainInfo),
	}
}

func (s *State) GetLastProposalsHeight(chain *types.Chain) int64 {
	if chainInfo, ok := s.ChainInfos[chain.Name]; !ok {
		return 0
	} else {
		return chainInfo.ProposalsHeight
	}
}

func (s *State) SetProposal(chain *types.Chain, proposal types.Proposal) {
	if _, ok := s.ChainInfos[chain.Name]; !ok {
		s.ChainInfos[chain.Name] = &ChainInfo{
			Chain:         chain,
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	s.ChainInfos[chain.Name].ProposalVotes[proposal.ID] = WalletVotes{
		Proposal: proposal,
		Votes:    make(map[string]ProposalVote),
	}
}

func (s *State) SetVote(chain *types.Chain, proposal types.Proposal, wallet *types.Wallet, vote ProposalVote) {
	if _, ok := s.ChainInfos[chain.Name]; !ok {
		s.ChainInfos[chain.Name] = &ChainInfo{
			Chain:         chain,
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	if _, ok := s.ChainInfos[chain.Name].ProposalVotes[proposal.ID]; !ok {
		s.ChainInfos[chain.Name].ProposalVotes[proposal.ID] = WalletVotes{
			Proposal: proposal,
			Votes:    make(map[string]ProposalVote),
		}
	}

	s.ChainInfos[chain.Name].ProposalVotes[proposal.ID].Votes[wallet.Address] = vote
}

func (s *State) SetChainProposalsError(chain *types.Chain, err *types.QueryError) {
	s.ChainInfos[chain.Name] = &ChainInfo{
		Chain:          chain,
		ProposalsError: err,
	}
}

func (s *State) SetChainProposalsHeight(
	chain *types.Chain,
	height int64,
) {
	if _, ok := s.ChainInfos[chain.Name]; !ok {
		s.ChainInfos[chain.Name] = &ChainInfo{
			Chain:         chain,
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	stateChain := s.ChainInfos[chain.Name]
	stateChain.ProposalsHeight = height
}

func (s *State) SetChainVotes(
	chain *types.Chain,
	votes map[string]WalletVotes,
) {
	stateChain := s.ChainInfos[chain.Name]
	stateChain.ProposalVotes = votes
}

func (s *State) GetProposal(chain, proposalID string) (types.Proposal, bool) {
	if _, ok := s.ChainInfos[chain]; !ok {
		return types.Proposal{}, false
	}
	chainInfo := s.ChainInfos[chain]

	if _, ok := chainInfo.ProposalVotes[proposalID]; !ok {
		return types.Proposal{}, false
	}

	return chainInfo.ProposalVotes[proposalID].Proposal, true
}

func (s *State) GetVote(chain, proposalID, wallet string) (ProposalVote, bool) {
	if _, ok := s.ChainInfos[chain]; !ok {
		return ProposalVote{}, false
	}
	chainInfo := s.ChainInfos[chain]

	if _, ok := chainInfo.ProposalVotes[proposalID]; !ok {
		return ProposalVote{}, false
	}
	proposalVotes := chainInfo.ProposalVotes[proposalID]

	vote, found := proposalVotes.Votes[wallet]
	return vote, found
}

func (s *State) HasVoted(chain, proposal, wallet string) bool {
	if _, ok := s.ChainInfos[chain]; !ok {
		return false
	}

	if _, ok := s.ChainInfos[chain].ProposalVotes[proposal]; !ok {
		return false
	}

	if _, ok := s.ChainInfos[chain].ProposalVotes[proposal].Votes[wallet]; !ok {
		return false
	}

	return s.ChainInfos[chain].ProposalVotes[proposal].Votes[wallet].HasVoted()
}

func (s *State) ToRenderedState() RenderedState {
	keys := make([]string, 0)
	renderedChainInfos := map[string]RenderedChainInfo{}

	for chainName, chainInfo := range s.ChainInfos {
		if !chainInfo.HasProposalsError() && len(chainInfo.ProposalVotes) == 0 {
			continue
		}

		proposalsKeys := make([]string, 0)
		renderedProposals := map[string]RenderedProposalVotes{}

		for proposalID, proposalVotes := range chainInfo.ProposalVotes {
			if !proposalVotes.Proposal.IsInVoting() {
				continue
			}

			votesKeys := make([]string, 0)
			renderedVotes := map[string]RenderedWalletVote{}

			for wallet, walletVote := range proposalVotes.Votes {
				votesKeys = append(votesKeys, wallet)
				renderedVotes[wallet] = RenderedWalletVote{
					Wallet: walletVote.Wallet,
					Vote:   walletVote.Vote,
					Error:  walletVote.Error,
				}
			}

			// sorting wallets votes by wallet name desc
			sort.Strings(votesKeys)

			proposalsKeys = append(proposalsKeys, proposalID)
			renderedProposals[proposalID] = RenderedProposalVotes{
				Proposal: proposalVotes.Proposal,
				Votes:    make([]RenderedWalletVote, len(votesKeys)),
			}

			for index, key := range votesKeys {
				renderedProposals[proposalID].Votes[index] = renderedVotes[key]
			}
		}

		keys = append(keys, chainName)
		renderedChainInfos[chainName] = RenderedChainInfo{
			Chain:          chainInfo.Chain,
			ProposalsError: chainInfo.ProposalsError,
			ProposalVotes:  make([]RenderedProposalVotes, len(proposalsKeys)),
		}

		// sorting proposals by ID desc
		sort.Slice(proposalsKeys, func(i, j int) bool {
			first, firstErr := strconv.Atoi(proposalsKeys[i])
			second, secondErr := strconv.Atoi(proposalsKeys[j])

			// if it's faulty - doesn't matter how we sort it out
			if firstErr != nil || secondErr != nil {
				return true
			}

			return first > second
		})

		for index, key := range proposalsKeys {
			renderedChainInfos[chainName].ProposalVotes[index] = renderedProposals[key]
		}
	}

	// sorting chains by chain name desc
	sort.Strings(keys)

	renderedState := RenderedState{
		ChainInfos: make([]RenderedChainInfo, len(keys)),
	}

	for index, key := range keys {
		renderedState.ChainInfos[index] = renderedChainInfos[key]
	}

	return renderedState
}
