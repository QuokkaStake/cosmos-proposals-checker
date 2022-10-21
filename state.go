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
	Proposal Proposal
	Votes    map[string]ProposalVote
}

type ChainInfo struct {
	Chain          Chain
	ProposalVotes  map[string]WalletVotes
	ProposalsError string
}

func (c ChainInfo) HasProposalsError() bool {
	return c.ProposalsError != ""
}

type State struct {
	ChainInfos map[string]ChainInfo
}

func NewState() State {
	return State{
		ChainInfos: make(map[string]ChainInfo),
	}
}

func (s *State) SetVote(chain Chain, proposal Proposal, wallet string, vote ProposalVote) {
	if _, ok := s.ChainInfos[chain.Name]; !ok {
		s.ChainInfos[chain.Name] = ChainInfo{
			Chain:         chain,
			ProposalVotes: make(map[string]WalletVotes),
		}
	}

	if _, ok := s.ChainInfos[chain.Name].ProposalVotes[proposal.ProposalID]; !ok {
		s.ChainInfos[chain.Name].ProposalVotes[proposal.ProposalID] = WalletVotes{
			Proposal: proposal,
			Votes:    make(map[string]ProposalVote),
		}
	}

	s.ChainInfos[chain.Name].ProposalVotes[proposal.ProposalID].Votes[wallet] = vote
}

func (s *State) SetChainProposalsError(chain Chain, err error) {
	s.ChainInfos[chain.Name] = ChainInfo{
		Chain:          chain,
		ProposalsError: err.Error(),
	}
}

func (s State) GetVoteAndProposal(chain, proposalID, wallet string) (ProposalVote, Proposal, bool) {
	if _, ok := s.ChainInfos[chain]; !ok {
		return ProposalVote{}, Proposal{}, false
	}

	if _, ok := s.ChainInfos[chain].ProposalVotes[proposalID]; !ok {
		return ProposalVote{}, Proposal{}, false
	}

	proposalVotes, found := s.ChainInfos[chain].ProposalVotes[proposalID]
	if !found {
		return ProposalVote{}, Proposal{}, false
	}

	vote, found := proposalVotes.Votes[wallet]
	return vote, proposalVotes.Proposal, found
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
