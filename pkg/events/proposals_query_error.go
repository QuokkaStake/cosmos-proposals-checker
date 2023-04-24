package events

import (
	configTypes "main/pkg/config/types"
)

type ProposalsQueryErrorEvent struct {
	Chain *configTypes.Chain
	Error error
}

func (e ProposalsQueryErrorEvent) Name() string {
	return "proposals_query_error"
}

func (e ProposalsQueryErrorEvent) IsAlert() bool {
	return false
}
