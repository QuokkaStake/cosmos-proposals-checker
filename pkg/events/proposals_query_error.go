package events

import (
	types "main/pkg/types"
)

type ProposalsQueryErrorEvent struct {
	Chain *types.Chain
	Error error
}

func (e ProposalsQueryErrorEvent) Name() string {
	return "proposals_query_error"
}

func (e ProposalsQueryErrorEvent) IsAlert() bool {
	return false
}
