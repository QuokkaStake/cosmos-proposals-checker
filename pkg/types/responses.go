package types

import (
	"main/pkg/utils"
	"time"
)

type Proposal struct {
	ProposalID    string           `json:"proposal_id"`
	Status        string           `json:"status"`
	Content       *ProposalContent `json:"content"`
	VotingEndTime time.Time        `json:"voting_end_time"`
}

type ProposalContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProposalsRPCResponse struct {
	Proposals []Proposal `json:"proposals"`
}

type Vote struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

func (v Vote) ResolveVote() string {
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
