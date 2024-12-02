package events

import (
	"main/pkg/types"
	"main/pkg/utils"
	"time"
)

type VotedEvent struct {
	RenderTime time.Time
	Chain      *types.Chain
	Wallet     *types.Wallet
	Proposal   types.Proposal
	Vote       *types.Vote
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

func (e VotedEvent) GetProposalTimeLeft() string {
	return utils.FormatDuration(e.Proposal.EndTime.Sub(e.RenderTime).Round(time.Second))
}
