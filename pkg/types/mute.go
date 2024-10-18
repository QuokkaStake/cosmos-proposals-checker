package types

import (
	"time"

	"github.com/guregu/null/v5"
)

type Mute struct {
	Chain      null.String
	ProposalID null.String
	Expires    time.Time
	Comment    string
}

func (m *Mute) IsExpired() bool {
	return m.Expires.Before(time.Now())
}

func (m *Mute) Matches(chain string, proposalID string) bool {
	match := true

	if !m.Chain.IsZero() {
		match = match && chain == m.Chain.String
	}

	if !m.ProposalID.IsZero() {
		match = match && proposalID == m.ProposalID.String
	}

	return match
}

func (m *Mute) LabelsEqual(another *Mute) bool {
	return m.Chain == another.Chain && m.ProposalID == another.ProposalID
}
