package state

import (
	"encoding/json"
	"main/pkg/fs"

	"github.com/rs/zerolog"
)

type Manager struct {
	Filesystem fs.FS
	StatePath  string
	Logger     zerolog.Logger
	State      State
}

func NewStateManager(path string, filesystem fs.FS, logger *zerolog.Logger) *Manager {
	return &Manager{
		StatePath:  path,
		Filesystem: filesystem,
		Logger:     logger.With().Str("component", "state_manager").Logger(),
		State:      NewState(),
	}
}

func (m *Manager) Load() {
	content, err := m.Filesystem.ReadFile(m.StatePath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load state")
		return
	}

	var s State
	if err = json.Unmarshal(content, &s); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall state")
		m.State = NewState()
		return
	}

	m.State = s
}

func (m *Manager) Save() {
	content, err := json.Marshal(m.State)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not marshal state")
		return
	}

	if err = m.Filesystem.WriteFile(m.StatePath, content, 0o600); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not save state")
		return
	}
}

func (m *Manager) CommitState(state State) {
	m.State = state
	m.Save()
}
