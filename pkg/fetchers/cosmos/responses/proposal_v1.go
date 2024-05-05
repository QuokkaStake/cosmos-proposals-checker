package responses

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
	"time"
)

// cosmos/gov/v1beta1/proposals?pagination.limit=1000&pagination.offset=0

type V1ProposalMessage struct {
	Content ProposalContent `json:"content"`
}

type V1Proposal struct {
	ProposalID    string              `json:"id"`
	Status        string              `json:"status"`
	VotingEndTime time.Time           `json:"voting_end_time"`
	Messages      []V1ProposalMessage `json:"messages"`

	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func (p V1Proposal) ToProposal() types.Proposal {
	// Some chains (namely, Quicksilver) do not have title and description fields,
	// instead they have content.title and content.description per each message.
	// Others (namely, Kujira) have title and summary text.
	// This should work for all of them.
	title := p.Title
	if title == "" {
		titles := utils.Map(p.Messages, func(m V1ProposalMessage) string {
			return m.Content.Title
		})

		title = strings.Join(titles, ", ")
	}

	description := p.Summary
	if description == "" {
		descriptions := utils.Map(p.Messages, func(m V1ProposalMessage) string {
			return m.Content.Description
		})

		description = strings.Join(descriptions, ", ")
	}

	return types.Proposal{
		ID:          p.ProposalID,
		Title:       title,
		Description: description,
		EndTime:     p.VotingEndTime,
		Status:      ParseProposalStatus(p.Status),
	}
}

type V1ProposalsRPCResponse struct {
	Code      int64        `json:"code"`
	Message   string       `json:"message"`
	Proposals []V1Proposal `json:"proposals"`
}
