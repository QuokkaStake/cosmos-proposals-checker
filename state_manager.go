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

func NewStateManager(path string, logger *zerolog.Logger) *StateManager {
	return &StateManager{
		StatePath: path,
		Logger:    logger.With().Str("component", "state_manager").Logger(),
		State:     NewState(),
	}
}

func (m *StateManager) Load() {
	content, err := os.ReadFile(m.StatePath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load state")
		return
	}

	var state State
	if err = json.Unmarshal(content, &state); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall state")
		m.State = NewState()
		return
	}

	m.State = state
}

func (m *StateManager) Save() {
	content, err := json.Marshal(m.State)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not marshal state")
		return
	}

	if err = os.WriteFile(m.StatePath, content, 0o600); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not save state")
		return
	}
}

func (m *StateManager) CommitState(state State) {
	m.State = state
	m.Save()
}
