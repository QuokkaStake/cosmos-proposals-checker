package entry

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
)

type ReportEntry interface {
	Name() string
	IsAlert() bool // only voted/not_voted are alerts, required for PagerDuty
}

type ReportEntryNotError interface {
	ReportEntry
	GetChain() *configTypes.Chain
	GetWallet() *configTypes.Wallet
	GetProposal() types.Proposal
}
