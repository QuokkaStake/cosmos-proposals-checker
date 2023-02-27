package mutesmanager

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
)

type Manager struct {
	MutesPath string
	Logger    zerolog.Logger
	Mutes     Mutes
}

func NewMutesManager(mutesPath string, logger *zerolog.Logger) *Manager {
	return &Manager{
		MutesPath: mutesPath,
		Logger:    logger.With().Str("component", "mutes_manager").Logger(),
		Mutes:     Mutes{},
	}
}

func (m *Manager) Load() {
	if m.MutesPath == "" {
		m.Logger.Debug().Msg("Mutes path not configured, not loading.")
		return
	}

	content, err := os.ReadFile(m.MutesPath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load mutes")
		return
	}

	var mutes Mutes
	if err = json.Unmarshal(content, &mutes); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall mutes")
		m.Mutes = Mutes{}
	}

	m.Mutes = mutes
}

func (m *Manager) Save() {
	if m.MutesPath == "" {
		m.Logger.Debug().Msg("Mutes path not configured, not saving.")
		return
	}

	content, err := json.Marshal(m.Mutes)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not marshal mutes")
		return
	}

	if err = os.WriteFile(m.MutesPath, content, 0o600); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not save mutes")
		return
	}
}

func (m *Manager) IsMuted(chain string, proposalID string) bool {
	if m.MutesPath == "" {
		return false
	}

	return m.Mutes.IsMuted(chain, proposalID)
}

func (m *Manager) AddMute(mute *Mute) {
	m.Mutes.AddMute(mute)
	m.Save()
}

func (m *Manager) DeleteMute(chain string, proposalID string) {
	m.Mutes.DeleteMute(chain, proposalID)
	m.Save()
}
