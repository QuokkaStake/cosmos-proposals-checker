package events

import (
	types "main/pkg/types"
)

type GenericErrorEvent struct {
	Chain *types.Chain
	Error error
}

func (e GenericErrorEvent) Name() string {
	return "generic_error"
}

func (e GenericErrorEvent) IsAlert() bool {
	return false
}
