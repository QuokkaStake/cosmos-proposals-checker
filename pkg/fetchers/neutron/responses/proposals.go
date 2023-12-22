package responses

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"
	"time"
)

type ProposalWithID struct {
	ID       int      `json:"id"`
	Proposal Proposal `json:"proposal"`
}
type Proposal struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Expiration  struct {
		AtTime string `json:"at_time"`
	} `json:"expiration"`
	Status     string `json:"status"`
	TotalPower string `json:"total_power"`
}

type ProposalsResponse struct {
	Data struct {
		Proposals []ProposalWithID `json:"proposals"`
	} `json:"data"`
}

func (p ProposalsResponse) ToProposals() ([]types.Proposal, error) {
	allProposals := utils.Filter(p.Data.Proposals, func(p ProposalWithID) bool {
		return p.Proposal.Status == "open"
	})

	proposals := make([]types.Proposal, len(allProposals))

	for index, proposal := range allProposals {
		expiresAt, err := strconv.ParseInt(proposal.Proposal.Expiration.AtTime, 10, 64)
		if err != nil {
			return nil, err
		}

		proposals[index] = types.Proposal{
			ID:          strconv.Itoa(proposal.ID),
			Title:       proposal.Proposal.Title,
			Description: proposal.Proposal.Description,
			EndTime:     time.Unix(0, expiresAt),
		}
	}

	return proposals, nil
}
