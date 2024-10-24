package mutesmanager

import (
	databasePkg "main/pkg/database"
	"main/pkg/report/entry"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Manager struct {
	Database databasePkg.Database
	Logger   zerolog.Logger
}

func NewMutesManager(logger *zerolog.Logger, database databasePkg.Database) *Manager {
	return &Manager{
		Database: database,
		Logger:   logger.With().Str("component", "mutes_manager").Logger(),
	}
}

func (m *Manager) IsEntryMuted(reportEntry entry.ReportEntry) (bool, error) {
	entryConverted, ok := reportEntry.(entry.ReportEntryNotError)
	if !ok {
		return false, nil
	}

	chain := entryConverted.GetChain()
	proposal := entryConverted.GetProposal()

	return m.Database.IsMuted(chain.Name, proposal.ID)
}

func (m *Manager) GetAllMutes() ([]*types.Mute, error) {
	return m.Database.GetAllMutes()
}

func (m *Manager) AddMute(mute *types.Mute) error {
	return m.Database.UpsertMute(mute)
}

func (m *Manager) DeleteMute(mute *types.Mute) (bool, error) {
	return m.Database.DeleteMute(mute)
}
