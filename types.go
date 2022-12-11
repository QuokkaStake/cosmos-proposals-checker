package main

import (
	"fmt"
	"time"
)

// RPC response types.
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
	return ResolveVote(v.Option)
}

type VoteRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Vote    *Vote  `json:"vote"`
}

func (v VoteRPCResponse) IsError() bool {
	return v.Code != 0
}

type Report struct {
	Entries []ReportEntry
}

func (r *Report) Empty() bool {
	return len(r.Entries) == 0
}

type ReportEntryType string

const (
	NotVoted           ReportEntryType = "not_voted"
	Voted              ReportEntryType = "voted"
	Revoted            ReportEntryType = "revoted"
	ProposalQueryError ReportEntryType = "proposal_query_error"
	VoteQueryError     ReportEntryType = "vote_query_error"
)

type ReportEntry struct {
	Chain                  Chain
	Wallet                 string
	ProposalID             string
	ProposalTitle          string
	ProposalDescription    string
	ProposalVoteEndingTime time.Time
	Type                   ReportEntryType
	Value                  string
	OldValue               string
}

func (e ReportEntry) HasVoted() bool {
	return e.Value != ""
}

func (e ReportEntry) HasRevoted() bool {
	return e.Value != "" && e.OldValue != ""
}

func (e ReportEntry) IsVoteOrNotVoted() bool {
	return e.Type == NotVoted || e.Type == Voted
}

func (e ReportEntry) GetProposalTime() string {
	return e.ProposalVoteEndingTime.Format(time.RFC3339Nano)
}

func (e ReportEntry) GetProposalTimeLeft() string {
	return time.Until(e.ProposalVoteEndingTime).Round(time.Second).String()
}

func (e ReportEntry) GetVote() string {
	return ResolveVote(e.Value)
}

func (e ReportEntry) GetOldVote() string {
	return ResolveVote(e.OldValue)
}

func ResolveVote(value string) string {
	votes := map[string]string{
		"VOTE_OPTION_YES":          "Yes",
		"VOTE_OPTION_ABSTAIN":      "Abstain",
		"VOTE_OPTION_NO":           "No",
		"VOTE_OPTION_NO_WITH_VETO": "No with veto",
	}

	if vote, ok := votes[value]; ok && vote != "" {
		return vote
	}

	return value
}

type Reporter interface {
	Init()
	Enabled() bool
	SendReport(report Report) error
	Name() string
}

type Link struct {
	Name string
	Href string
}

func (l Link) Serialize() string {
	if l.Href == "" {
		return l.Name
	}

	return fmt.Sprintf("<a href='%s'>%s</a>", l.Href, l.Name)
}

type Mute struct {
	Chain      string
	ProposalID string
	Expires    time.Time
	Comment    string
}

type Mutes struct {
	Mutes []Mute
}

func (m *Mute) IsExpired() bool {
	return m.Expires.Before(time.Now())
}

func (m *Mute) Matches(chain string, proposalID string) bool {
	if m.Chain != chain {
		return false
	}

	// whole chain is muted
	if m.ProposalID == "" {
		return true
	}

	return m.ProposalID == proposalID
}

func (m *Mutes) IsMuted(chain string, proposalID string) bool {
	for _, mute := range m.Mutes {
		if mute.Matches(chain, proposalID) {
			return !mute.IsExpired()
		}
	}

	return false
}

func (m Mute) GetExpirationTime() string {
	return m.Expires.Format(time.RFC3339Nano)
}

func (m *Mutes) AddMute(mute Mute) {
	m.Mutes = append(m.Mutes, mute)
	m.Mutes = Filter(m.Mutes, func(m Mute) bool {
		return !m.IsExpired()
	})
}

func (m *Mutes) DeleteMute(chain string, proposalID string) bool {
	for index, mute := range m.Mutes {
		if mute.Chain == chain && mute.ProposalID == proposalID {
			m.Mutes = append(m.Mutes[:index], m.Mutes[index+1:]...)
			m.Mutes = Filter(m.Mutes, func(m Mute) bool {
				return !m.IsExpired()
			})
			return true
		}
	}

	m.Mutes = Filter(m.Mutes, func(m Mute) bool {
		return !m.IsExpired()
	})

	return false
}
