package database

import (
	"database/sql"
	"errors"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/rs/zerolog"
)

type Logger struct {
	Logger zerolog.Logger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(strings.TrimSpace(format), v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Panic().Msgf(strings.TrimSpace(format), v...)
}

type Database struct {
	logger         zerolog.Logger
	config         types.DatabaseConfig
	client         *sql.DB
	databaseLogger goose.Logger
}

func NewDatabase(
	logger *zerolog.Logger,
	config types.DatabaseConfig,
) *Database {
	return &Database{
		logger: logger.With().Str("component", "database").Logger(),
		config: config,
		databaseLogger: &Logger{
			Logger: logger.With().Str("component", "migrations").Logger(),
		},
	}
}

func (d *Database) Init() {
	db, err := sql.Open("sqlite3", d.config.Path)

	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not open SQLite database")
	}

	var version string
	if versionErr := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version); versionErr != nil {
		d.logger.Panic().Err(err).Msg("Could not query SQLite database")
	}

	d.logger.Info().
		Str("version", version).
		Str("path", d.config.Path).
		Msg("SQLite database connected")

	d.client = db
}

func (d *Database) Migrate() {
	goose.SetBaseFS(migrationsPkg.EmbedFS)
	goose.SetLogger(d.databaseLogger)

	_ = goose.SetDialect("sqlite3")

	if err := goose.Up(d.client, "."); err != nil {
		d.logger.Panic().Err(err).Msg("Could not apply migrations")
	}
}

func (d *Database) Rollback() {
	goose.SetBaseFS(migrationsPkg.EmbedFS)
	goose.SetLogger(d.databaseLogger)

	_ = goose.SetDialect("sqlite3")

	if err := goose.Reset(d.client, "."); err != nil {
		if errors.Is(err, goose.ErrNoCurrentVersion) {
			d.logger.Info().Err(err).Msg("No migrations are applied, cannot rollback")
		} else {
			d.logger.Panic().Err(err).Msg("Could not rollback migrations")
		}
	}
}

func (d *Database) UpsertProposal(
	chain *types.Chain,
	proposal types.Proposal,
) error {
	_, err := d.client.Exec(
		"INSERT INTO proposals (chain, id, title, description, status, end_time) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO UPDATE SET title = $3, description = $4, status = $5, end_time = $6",
		chain.Name,
		proposal.ID,
		proposal.Title,
		proposal.Description,
		proposal.Status,
		proposal.EndTime,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not upsert proposal")
		return err
	}

	return nil
}

func (d *Database) GetProposal(chain *types.Chain, proposalID string) (*types.Proposal, error) {
	proposal := &types.Proposal{}
	row := d.client.QueryRow(
		"SELECT id, title, description, status, end_time FROM proposals WHERE chain = $1 AND id = $2 LIMIT 1",
		chain.Name,
		proposalID,
	)

	err := row.Scan(
		&proposal.ID,
		&proposal.Title,
		&proposal.Description,
		&proposal.Status,
		&proposal.EndTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		d.logger.Error().Err(err).Msg("Error getting proposal")
		return nil, err
	}

	return proposal, nil
}
