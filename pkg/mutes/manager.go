package mutesmanager

import (
	"encoding/json"
	"main/pkg/fs"
	"main/pkg/utils"

	"github.com/rs/zerolog"
)

type Manager struct {
	Filesystem fs.FS
	MutesPath  string
	Logger     zerolog.Logger
	Mutes      Mutes
}

func NewMutesManager(mutesPath string, filesystem fs.FS, logger *zerolog.Logger) *Manager {
	return &Manager{
		Filesystem: filesystem,
		MutesPath:  mutesPath,
		Logger:     logger.With().Str("component", "mutes_manager").Logger(),
		Mutes:      Mutes{},
	}
}

func (m *Manager) Load() {
	if m.MutesPath == "" {
		m.Logger.Debug().Msg("Mutes path not configured, not loading.")
		return
	}

	content, err := m.Filesystem.ReadFile(m.MutesPath)
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

	content := utils.MustMarshal(m.Mutes)

	if err := m.Filesystem.WriteFile(m.MutesPath, content, 0o600); err != nil {
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
