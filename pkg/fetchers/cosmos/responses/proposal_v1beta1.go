package responses

import (
	"main/pkg/types"
	"time"
)

// cosmos/gov/v1beta1/proposals?pagination.limit=1000&pagination.offset=0

type V1beta1Proposal struct {
	ProposalID    string           `json:"proposal_id"`
	Status        string           `json:"status"`
	Content       *ProposalContent `json:"content"`
	VotingEndTime time.Time        `json:"voting_end_time"`
}

type ProposalContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type V1Beta1ProposalsRPCResponse struct {
	Code      int64             `json:"code"`
	Message   string            `json:"message"`
	Proposals []V1beta1Proposal `json:"proposals"`
}

func (p V1beta1Proposal) ToProposal() types.Proposal {
	return types.Proposal{
		ID:          p.ProposalID,
		Title:       p.Content.Title,
		Description: p.Content.Description,
		EndTime:     p.VotingEndTime,
	}
}
