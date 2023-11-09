package events

import (
	"main/pkg/types"
)

type RevotedEvent struct {
	Chain    *types.Chain
	Wallet   *types.Wallet
	Proposal types.Proposal
	Vote     *types.Vote
	OldVote  *types.Vote
}

func (e RevotedEvent) Name() string {
	return "revoted"
}

func (e RevotedEvent) IsAlert() bool {
	return true
}

func (e RevotedEvent) GetChain() *types.Chain {
	return e.Chain
}

func (e RevotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e RevotedEvent) GetWallet() *types.Wallet {
	return e.Wallet
}
