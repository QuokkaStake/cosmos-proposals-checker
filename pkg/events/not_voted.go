package events

import (
	"main/pkg/types"
)

type NotVotedEvent struct {
	Chain    *types.Chain
	Wallet   *types.Wallet
	Proposal types.Proposal
}

func (e NotVotedEvent) Name() string {
	return "not_voted"
}

func (e NotVotedEvent) IsAlert() bool {
	return true
}

func (e NotVotedEvent) GetChain() *types.Chain {
	return e.Chain
}

func (e NotVotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e NotVotedEvent) GetWallet() *types.Wallet {
	return e.Wallet
}
