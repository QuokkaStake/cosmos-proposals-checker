package main

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

type Report struct {
	Entries []ReportEntry
}

func (r *Report) Empty() bool {
	return len(r.Entries) == 0
}

type ReportEntry struct {
	Chain               Chain
	Wallet              string
	ProposalID          string
	ProposalTitle       string
	ProposalDescription string
	Vote                string
}

func (e *ReportEntry) HasVoted() bool {
	return e.Vote != ""
}

type Reporter interface {
	Init()
	Enabled() bool
	SendReport(report Report) error
	Name() string
}
