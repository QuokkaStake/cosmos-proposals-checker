package types

import (
	"main/pkg/utils"
	"time"
)

type ProposalStatus string

const (
	ProposalStatusVoting   ProposalStatus = "voting"
	ProposalStatusDeposit  ProposalStatus = "deposit"
	ProposalStatusPassed   ProposalStatus = "passed"
	ProposalStatusRejected ProposalStatus = "rejected"
	ProposalStatusFailed   ProposalStatus = "failed"
)

type Proposal struct {
	ID          string
	Title       string
	Description string
	EndTime     time.Time
	Status      ProposalStatus
}

func (p Proposal) GetTimeLeft() string {
	return utils.FormatDuration(time.Until(p.EndTime).Round(time.Second))
}

func (p ProposalStatus) String() string {
	switch p {
	case ProposalStatusVoting:
		return "🗳️Voting"
	case ProposalStatusDeposit:
		return "💸Deposit"
	case ProposalStatusPassed:
		return "🙌 Passed"
	case ProposalStatusRejected:
		return "🙅‍Rejected"
	case ProposalStatusFailed:
		return "🤦‍Failed"
	default:
		return string(p)
	}
}

func (p Proposal) IsInVoting() bool {
	return p.Status == ProposalStatusVoting
}
