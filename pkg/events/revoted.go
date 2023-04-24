package events

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
)

type RevotedEvent struct {
	Chain    *configTypes.Chain
	Wallet   *configTypes.Wallet
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

func (e RevotedEvent) GetChain() *configTypes.Chain {
	return e.Chain
}

func (e RevotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e RevotedEvent) GetWallet() *configTypes.Wallet {
	return e.Wallet
}
