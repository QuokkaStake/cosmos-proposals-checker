package events

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
)

type VotedEvent struct {
	Chain    *configTypes.Chain
	Wallet   *configTypes.Wallet
	Proposal types.Proposal
	Vote     *types.Vote
}

func (e VotedEvent) Name() string {
	return "voted"
}

func (e VotedEvent) IsAlert() bool {
	return true
}

func (e VotedEvent) GetChain() *configTypes.Chain {
	return e.Chain
}

func (e VotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e VotedEvent) GetWallet() *configTypes.Wallet {
	return e.Wallet
}
