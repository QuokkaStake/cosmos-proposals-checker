package events

import (
	"main/pkg/types"
	"main/pkg/utils"
	"time"
)

type RevotedEvent struct {
	RenderTime time.Time
	Chain      *types.Chain
	Wallet     *types.Wallet
	Proposal   types.Proposal
	Vote       *types.Vote
	OldVote    *types.Vote
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

func (e NotVotedEvent) GetProposalTimeLeft() string {
	return utils.FormatDuration(e.Proposal.EndTime.Sub(e.RenderTime).Round(time.Second))
}
