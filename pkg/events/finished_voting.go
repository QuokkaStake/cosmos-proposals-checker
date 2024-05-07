package events

import (
	"main/pkg/types"
)

type FinishedVotingEvent struct {
	Chain    *types.Chain
	Proposal types.Proposal
}

func (e FinishedVotingEvent) Name() string {
	return "finished_voting"
}

func (e FinishedVotingEvent) IsAlert() bool {
	return false
}
