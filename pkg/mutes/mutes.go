package mutesmanager

import (
	"github.com/guregu/null/v5"
	"time"

	"main/pkg/utils"
)

type Mute struct {
	Chain      null.String
	ProposalID null.String
	Expires    time.Time
	Comment    string
}

type Mutes struct {
	Mutes []*Mute
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

func (m *Mutes) IsMuted(chain string, proposalID string) bool {
	for _, mute := range m.Mutes {
		if mute.IsExpired() {
			continue
		}

		if mute.Matches(chain, proposalID) {
			return true
		}
	}

	return false
}

func (m *Mutes) AddMute(mute *Mute) {
	for _, existingMute := range m.Mutes {
		if existingMute.LabelsEqual(mute) {
			existingMute.Comment = mute.Comment
			existingMute.Expires = mute.Expires

			m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
				return !m.IsExpired()
			})
			return
		}
	}

	m.Mutes = append(m.Mutes, mute)
	m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
		return !m.IsExpired()
	})
}

func (m *Mutes) DeleteMute(mute *Mute) bool {
	for index, existingMute := range m.Mutes {
		if existingMute.LabelsEqual(mute) {
			m.Mutes = append(m.Mutes[:index], m.Mutes[index+1:]...)
			m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
				return !m.IsExpired()
			})
			return true
		}
	}

	m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
		return !m.IsExpired()
	})
	return false
}
