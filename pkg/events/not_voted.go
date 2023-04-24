package events

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
)

type NotVotedEvent struct {
	Chain    *configTypes.Chain
	Wallet   *configTypes.Wallet
	Proposal types.Proposal
}

func (e NotVotedEvent) Name() string {
	return "not_voted"
}

func (e NotVotedEvent) IsAlert() bool {
	return true
}

func (e NotVotedEvent) GetChain() *configTypes.Chain {
	return e.Chain
}

func (e NotVotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e NotVotedEvent) GetWallet() *configTypes.Wallet {
	return e.Wallet
}
