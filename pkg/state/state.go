package state

import (
	"main/pkg/types"
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

func (s *State) GetVoteAndProposal(chain, proposalID, wallet string) (ProposalVote, types.Proposal, bool) {
	if _, ok := s.ChainInfos[chain]; !ok {
		return ProposalVote{}, types.Proposal{}, false
	}
	chainInfo := s.ChainInfos[chain]

	if _, ok := chainInfo.ProposalVotes[proposalID]; !ok {
		return ProposalVote{}, types.Proposal{}, false
	}
	proposalVotes := chainInfo.ProposalVotes[proposalID]

	vote, found := proposalVotes.Votes[wallet]
	return vote, proposalVotes.Proposal, found
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
