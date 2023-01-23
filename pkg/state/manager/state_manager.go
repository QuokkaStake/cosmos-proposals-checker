package manager

import (
	"encoding/json"
	"os"

	"main/pkg/state"

	"github.com/rs/zerolog"
)

type StateManager struct {
	StatePath string
	Logger    zerolog.Logger
	State     state.State
}

func NewStateManager(path string, logger *zerolog.Logger) *StateManager {
	return &StateManager{
		StatePath: path,
		Logger:    logger.With().Str("component", "state_manager").Logger(),
		State:     state.NewState(),
	}
}

func (m *StateManager) Load() {
	content, err := os.ReadFile(m.StatePath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load state")
		return
	}

	var s state.State
	if err = json.Unmarshal(content, &s); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall state")
		m.State = state.NewState()
		return
	}

	m.State = s
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

func (m *StateManager) CommitState(state state.State) {
	m.State = state
	m.Save()
}
