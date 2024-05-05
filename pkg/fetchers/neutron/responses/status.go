package responses

import "main/pkg/types"

func ParseProposalStatus(status string) types.ProposalStatus {
	switch status {
	case "PROPOSAL_STATUS_DEPOSIT_PERIOD":
		return types.ProposalStatusDeposit
	case "open":
		return types.ProposalStatusVoting
	case "executed":
		return types.ProposalStatusPassed
	case "rejected":
		return types.ProposalStatusRejected
	case "execution_failed":
		return types.ProposalStatusFailed
	default:
		return types.ProposalStatus(status)
	}
}
