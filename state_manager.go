package main

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
)

type StateManager struct {
	StatePath string
	Logger    zerolog.Logger
	State     State
}

type State struct {
	VotesState    VotesState
	OldVotesState VotesState
}

type (
	WalletVotes   map[string]*Vote
	ProposalVotes map[string]WalletVotes
)

// ["chain"]["proposal"]["wallet"]["vote"].
type VotesState map[string]ProposalVotes

func NewStateManager(path string, logger *zerolog.Logger) *StateManager {
	return &StateManager{
		StatePath: path,
		Logger:    logger.With().Str("component", "state_manager").Logger(),
		State: State{
			VotesState:    make(VotesState),
			OldVotesState: make(VotesState),
		},
	}
}

func (m *StateManager) SetVote(chain, proposal, wallet string, vote *Vote) {
	votesState := m.State.VotesState

	if _, ok := votesState[chain]; !ok {
		votesState[chain] = make(ProposalVotes)
	}

	if _, ok := votesState[chain][proposal]; !ok {
		votesState[chain][proposal] = make(WalletVotes)
	}

	if vote != nil {
		votesState[chain][proposal][wallet] = vote
	}
}

func (m *StateManager) GetVote(chain, proposal, wallet string) *Vote {
	votesState := m.State.VotesState

	if _, ok := votesState[chain]; !ok {
		return nil
	}

	if _, ok := votesState[chain][proposal]; !ok {
		return nil
	}

	return votesState[chain][proposal][wallet]
}

func (m *StateManager) HasVotedNow(chain, proposal, wallet string) bool {
	if m.State.VotesState == nil {
		return false
	}

	votesState := m.State.VotesState
	if _, ok := votesState[chain]; !ok {
		return false
	}

	if _, ok := votesState[chain][proposal]; !ok {
		return false
	}

	_, ok := votesState[chain][proposal][wallet]
	return ok
}

func (m *StateManager) HasVotedBefore(chain, proposal, wallet string) bool {
	if m.State.OldVotesState == nil {
		return false
	}

	votesState := m.State.OldVotesState
	if _, ok := votesState[chain]; !ok {
		return false
	}

	if _, ok := votesState[chain][proposal]; !ok {
		return false
	}

	_, ok := votesState[chain][proposal][wallet]
	return ok
}

func (m *StateManager) Load() {
	content, err := os.ReadFile(m.StatePath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load state")
		return
	}

	var state VotesState
	if err = json.Unmarshal(content, &state); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall state")
		m.State.OldVotesState = make(VotesState)
		return
	}

	m.State.OldVotesState = state
	m.State.VotesState = make(VotesState)
}

func (m *StateManager) Save() {
	content, err := json.Marshal(m.State.OldVotesState)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not marshal state")
		return
	}

	if err = os.WriteFile(m.StatePath, content, 0o600); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not save state")
		return
	}
}

func (m *StateManager) CommitNewState() {
	m.State.OldVotesState = m.State.VotesState
	m.State.VotesState = make(VotesState)

	m.Save()
}
