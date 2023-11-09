package events

import (
	"main/pkg/types"
)

type VotedEvent struct {
	Chain    *types.Chain
	Wallet   *types.Wallet
	Proposal types.Proposal
	Vote     *types.Vote
}

func (e VotedEvent) Name() string {
	return "voted"
}

func (e VotedEvent) IsAlert() bool {
	return true
}

func (e VotedEvent) GetChain() *types.Chain {
	return e.Chain
}

func (e VotedEvent) GetProposal() types.Proposal {
	return e.Proposal
}

func (e VotedEvent) GetWallet() *types.Wallet {
	return e.Wallet
}
