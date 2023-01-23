package mutesmanager

import (
	"encoding/json"
	"os"

	"main/pkg/types"

	"github.com/rs/zerolog"
)

type MutesManager struct {
	MutesPath string
	Logger    zerolog.Logger
	Mutes     types.Mutes
}

func NewMutesManager(mutesPath string, logger *zerolog.Logger) *MutesManager {
	return &MutesManager{
		MutesPath: mutesPath,
		Logger:    logger.With().Str("component", "mutes_manager").Logger(),
		Mutes:     types.Mutes{},
	}
}

func (m *MutesManager) Load() {
	if m.MutesPath == "" {
		m.Logger.Debug().Msg("Mutes path not configured, not loading.")
		return
	}

	content, err := os.ReadFile(m.MutesPath)
	if err != nil {
		m.Logger.Warn().Err(err).Msg("Could not load mutes")
		return
	}

	var mutes types.Mutes
	if err = json.Unmarshal(content, &mutes); err != nil {
		m.Logger.Warn().Err(err).Msg("Could not unmarshall mutes")
		m.Mutes = types.Mutes{}
	}

	m.Mutes = mutes
}

func (m *MutesManager) Save() {
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

func (m *MutesManager) IsMuted(chain string, proposalID string) bool {
	if m.MutesPath == "" {
		return false
	}

	return m.Mutes.IsMuted(chain, proposalID)
}

func (m *MutesManager) AddMute(mute types.Mute) {
	m.Mutes.AddMute(mute)
	m.Save()
}

func (m *MutesManager) DeleteMute(chain string, proposalID string) {
	m.Mutes.DeleteMute(chain, proposalID)
	m.Save()
}
