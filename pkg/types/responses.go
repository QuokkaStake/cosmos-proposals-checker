package types

import (
	"strings"
	"time"

	"main/pkg/utils"
)

type Proposal struct {
	ProposalID    string           `json:"proposal_id"`
	Status        string           `json:"status"`
	Content       *ProposalContent `json:"content"`
	VotingEndTime time.Time        `json:"voting_end_time"`
}

func (p Proposal) GetTimeLeft() string {
	return utils.FormatDuration(time.Until(p.VotingEndTime).Round(time.Second))
}

func (p Proposal) GetProposalTime() string {
	return p.VotingEndTime.Format(time.RFC1123)
}

type ProposalContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProposalsRPCResponse struct {
	Code      int64      `json:"code"`
	Message   string     `json:"message"`
	Proposals []Proposal `json:"proposals"`
}

type Vote struct {
	ProposalID string       `json:"proposal_id"`
	Voter      string       `json:"voter"`
	Option     string       `json:"option"`
	Options    []VoteOption `json:"options"`
}

type VoteOption struct {
	Option string `json:"option"`
	Weight string `json:"weight"`
}

func (v Vote) ResolveVote() string {
	if len(v.Options) > 0 {
		optionsStrings := utils.Map(v.Options, func(v VoteOption) string {
			return utils.ResolveVote(v.Option)
		})

		return strings.Join(optionsStrings, ", ")
	}

	return utils.ResolveVote(v.Option)
}

type VoteRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Vote    *Vote  `json:"vote"`
}

func (v VoteRPCResponse) IsError() bool {
	return v.Code != 0
}
