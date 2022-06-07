package main

const VotingPeriod = "PROPOSAL_STATUS_VOTING_PERIOD"

// RPC response types.
type Proposal struct {
	ProposalID string           `json:"proposal_id"`
	Status     string           `json:"status"`
	Content    *ProposalContent `json:"content"`
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

type VoteRPCResponse struct {
	Code int64 `json:"code"`
	Vote *Vote `json:"vote"`
}
