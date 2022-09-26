package main

type ProposalVote struct {
	Vote  *Vote
	Error string
}

func (v ProposalVote) HasVoted() bool {
	return v.Vote != nil && v.Error == ""
}

func (v ProposalVote) IsError() bool {
	return v.Error != ""
}

type WalletVotes struct {
	Votes map[string]ProposalVote
}

type ChainInfo struct {
	Proposals      map[string]Proposal
	ProposalVotes  map[string]WalletVotes
	ProposalsError string
}

func (c *ChainInfo) HasProposalsError() bool {
	return c.ProposalsError != ""
}

type State struct {
	ChainInfos map[string]*ChainInfo
}

func NewState() State {
	return State{
		ChainInfos: make(map[string]*ChainInfo),
	}
}

func (s *State) SetVote(chain string, proposal Proposal, wallet string, vote ProposalVote) {
	if _, ok := s.ChainInfos[chain]; !ok {
		s.ChainInfos[chain] = &ChainInfo{
			Proposals:     make(map[string]Proposal),
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	s.ChainInfos[chain].Proposals[proposal.ProposalID] = proposal

	if _, ok := s.ChainInfos[chain].ProposalVotes[proposal.ProposalID]; !ok {
		s.ChainInfos[chain].ProposalVotes[proposal.ProposalID] = WalletVotes{
			Votes: make(map[string]ProposalVote),
		}
	}

	s.ChainInfos[chain].ProposalVotes[proposal.ProposalID].Votes[wallet] = vote
}

func (s *State) SetChainProposalsError(chain string, err error) {
	if _, ok := s.ChainInfos[chain]; !ok {
		s.ChainInfos[chain] = &ChainInfo{
			Proposals:     make(map[string]Proposal),
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	chainInfo := s.ChainInfos[chain]
	chainInfo.ProposalsError = err.Error()
}

func (s State) GetVoteAndProposal(chain, proposalId, wallet string) (ProposalVote, Proposal, bool) {
	if _, ok := s.ChainInfos[chain]; !ok {
		return ProposalVote{}, Proposal{}, false
	}

	if _, ok := s.ChainInfos[chain].ProposalVotes[proposalId]; !ok {
		return ProposalVote{}, Proposal{}, false
	}

	vote, found := s.ChainInfos[chain].ProposalVotes[proposalId].Votes[wallet]
	proposal := s.ChainInfos[chain].Proposals[proposalId]
	return vote, proposal, found
}

func (s State) HasVoted(chain, proposal, wallet string) bool {
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
