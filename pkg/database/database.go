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
