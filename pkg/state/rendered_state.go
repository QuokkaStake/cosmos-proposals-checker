package state

import "main/pkg/types"

type RenderedState struct {
	ChainInfos []RenderedChainInfo
}

type RenderedChainInfo struct {
	Chain          *types.Chain
	ProposalVotes  []RenderedProposalVotes
	ProposalsError *types.QueryError
}

type RenderedProposalVotes struct {
	Proposal types.Proposal
	Votes    []RenderedWalletVote
}

type RenderedWalletVote struct {
	Wallet *types.Wallet
	Vote   *types.Vote
	Error  *types.QueryError
}

func (v RenderedWalletVote) HasVoted() bool {
	return v.Vote != nil && v.Error == nil
}

func (v RenderedWalletVote) IsError() bool {
	return v.Error != nil
}

func (c RenderedChainInfo) HasProposalsError() bool {
	return c.ProposalsError != nil
}
