package responses

import "main/pkg/types"

func ParseProposalStatus(status string) types.ProposalStatus {
	switch status {
	case "PROPOSAL_STATUS_DEPOSIT_PERIOD":
		return types.ProposalStatusDeposit
	case "PROPOSAL_STATUS_VOTING_PERIOD":
		return types.ProposalStatusVoting
	case "PROPOSAL_STATUS_PASSED":
		return types.ProposalStatusPassed
	case "PROPOSAL_STATUS_REJECTED":
		return types.ProposalStatusRejected
	case "PROPOSAL_STATUS_FAILED":
		return types.ProposalStatusFailed
	default:
		return types.ProposalStatus(status)
	}
}
