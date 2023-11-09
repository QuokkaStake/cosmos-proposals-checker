package entry

import (
	"main/pkg/types"
)

type ReportEntry interface {
	Name() string
	IsAlert() bool // only voted/not_voted are alerts, required for PagerDuty
}

type ReportEntryNotError interface {
	ReportEntry
	GetChain() *types.Chain
	GetWallet() *types.Wallet
	GetProposal() types.Proposal
}
