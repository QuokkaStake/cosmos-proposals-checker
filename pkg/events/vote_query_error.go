package events

import (
	"main/pkg/types"
)

type VoteQueryError struct {
	Chain    *types.Chain
	Proposal types.Proposal
	Error    error
}

func (e VoteQueryError) Name() string {
	return "vote_query_error"
}

func (e VoteQueryError) IsAlert() bool {
	return true
}
