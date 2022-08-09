package main

import (
	"time"
)

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

type ExplorerLink struct {
	Name string
	Link string
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

func (m *Mutes) IsMuted(chain string, proposalID string) bool {
	for _, mute := range m.Mutes {
		if mute.Chain == chain && mute.ProposalID == proposalID {
			return !mute.IsExpired()
		}
	}

	return false
}

func (m *Mutes) AddMute(mute Mute) {
	for index, muteInRange := range m.Mutes {
		if mute.Chain == muteInRange.Chain && mute.ProposalID == muteInRange.ProposalID {
			m.Mutes[index] = mute
			return
		}
	}

	m.Mutes = append(m.Mutes, mute)
}

func (m *Mutes) DeleteMute(chain string, proposalID string) bool {
	for index, mute := range m.Mutes {
		if mute.Chain == chain && mute.ProposalID == proposalID {
			m.Mutes = append(m.Mutes[:index], m.Mutes[index+1:]...)
			return true
		}
	}

	return false
}
