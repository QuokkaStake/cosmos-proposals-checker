package database

import (
	"context"
	"database/sql"
	"errors"
	migrationsPkg "main/migrations"
	"main/pkg/types"
	"os"
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

type SqliteDatabase struct {
	logger         zerolog.Logger
	config         types.DatabaseConfig
	client         *sql.DB
	databaseLogger goose.Logger
}

func NewSqliteDatabase(
	logger *zerolog.Logger,
	config types.DatabaseConfig,
) *SqliteDatabase {
	return &SqliteDatabase{
		logger: logger.With().Str("component", "database").Logger(),
		config: config,
		databaseLogger: &Logger{
			Logger: logger.With().Str("component", "migrations").Logger(),
		},
	}
}

func (d *SqliteDatabase) Init() {
	db, err := sql.Open("sqlite3", d.config.Path)

	if err != nil {
		d.logger.Panic().Err(err).Msg("Could not open SQLite database")
	}

	var version string
	if versionErr := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version); versionErr != nil {
		d.logger.Panic().Err(versionErr).Msg("Could not query SQLite database")
	}

	d.logger.Info().
		Str("version", version).
		Str("path", d.config.Path).
		Msg("SQLite database connected")

	d.client = db
}

func (d *SqliteDatabase) Migrate() {
	goose.SetBaseFS(migrationsPkg.EmbedFS)
	goose.SetLogger(d.databaseLogger)

	_ = goose.SetDialect("sqlite3")

	if err := goose.Up(d.client, "."); err != nil {
		d.logger.Panic().Err(err).Msg("Could not apply migrations")
	}
}

func (d *SqliteDatabase) Rollback() {
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

func (d *SqliteDatabase) UpsertProposal(
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

func (d *SqliteDatabase) GetProposal(chain *types.Chain, proposalID string) (*types.Proposal, error) {
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

func (d *SqliteDatabase) GetVote(
	chain *types.Chain,
	proposal types.Proposal,
	wallet *types.Wallet,
) (*types.Vote, error) {
	rows, err := d.client.Query(
		"SELECT vote_option, vote_weight FROM votes WHERE chain = $1 AND proposal_id = $2 AND wallet = $3",
		chain.Name,
		proposal.ID,
		wallet.Address,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting vote")
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	vote := &types.Vote{
		ProposalID: proposal.ID,
		Voter:      wallet.Address,
		Options:    make(types.VoteOptions, 0),
	}

	for rows.Next() {
		voteOption := types.VoteOption{}

		scanErr := rows.Scan(&voteOption.Option, &voteOption.Weight)
		if scanErr != nil {
			d.logger.Error().Err(scanErr).Msg("Error getting vote")
			return nil, scanErr
		}

		vote.Options = append(vote.Options, voteOption)
	}

	if len(vote.Options) == 0 {
		return nil, nil //nolint:nilnil
	}

	return vote, nil
}

func (d *SqliteDatabase) UpsertVote(
	chain *types.Chain,
	proposal types.Proposal,
	wallet *types.Wallet,
	vote *types.Vote,
	ctx context.Context,
) error {
	tx, err := d.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint:errcheck

	if _, deleteErr := tx.Exec(
		"DELETE FROM votes WHERE chain = $1 AND proposal_id = $2 AND wallet = $3",
		chain.Name,
		proposal.ID,
		wallet.Address,
	); deleteErr != nil {
		d.logger.Error().Err(deleteErr).Msg("Error deleting votes")
		return err
	}

	for _, option := range vote.Options {
		if _, insertErr := tx.Exec(
			"INSERT INTO votes (chain, proposal_id, wallet, vote_option, vote_weight) VALUES ($1, $2, $3, $4, $5)",
			chain.Name,
			proposal.ID,
			wallet.Address,
			option.Option,
			option.Weight,
		); insertErr != nil {
			d.logger.Error().Err(insertErr).Msg("Error inserting votes")
			return err
		}
	}

	if insertErr := tx.Commit(); insertErr != nil {
		d.logger.Error().Err(insertErr).Msg("Error committing votes")
	}

	return nil
}

func (d *SqliteDatabase) GetLastBlockHeight(
	chain *types.Chain,
	storableKey string,
) (int64, error) {
	row := d.client.QueryRow(
		"SELECT height FROM query_last_block WHERE chain = $1 AND query = $2 LIMIT 1",
		chain.Name,
		storableKey,
	)

	var lastBlockHeight int64 = 0

	if err := row.Scan(&lastBlockHeight); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		d.logger.Error().Err(err).Msg("Error getting last block height")
		return 0, err
	}

	return lastBlockHeight, nil
}

func (d *SqliteDatabase) UpsertLastBlockHeight(
	chain *types.Chain,
	storableKey string,
	height int64,
) error {
	_, err := d.client.Exec(
		"INSERT INTO query_last_block (chain, query, height) VALUES ($1, $2, $3) ON CONFLICT DO UPDATE SET height = $3",
		chain.Name,
		storableKey,
		height,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not upsert last block height")
		return err
	}

	return nil
}

func (d *SqliteDatabase) UpsertMute(mute *types.Mute) error {
	_, err := d.client.Exec(
		"INSERT INTO mutes (chain, proposal_id, expires, comment) VALUES ($1, $2, $3, $4) ON CONFLICT DO UPDATE SET expires = $3, comment = $4",
		mute.Chain,
		mute.ProposalID,
		mute.Expires,
		mute.Comment,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not upsert mute")
		return err
	}

	return nil
}

func (d *SqliteDatabase) GetAllMutes() ([]*types.Mute, error) {
	mutes := make([]*types.Mute, 0)

	rows, err := d.client.Query("SELECT chain, proposal_id, expires, comment FROM mutes")
	if err != nil {
		d.logger.Error().Err(err).Msg("Error getting all mutes")
		return mutes, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		mute := &types.Mute{}

		err = rows.Scan(&mute.Chain, &mute.ProposalID, &mute.Expires, &mute.Comment)
		if err != nil {
			d.logger.Error().Err(err).Msg("Error getting mute")
			return mutes, err
		}

		mutes = append(mutes, mute)
	}

	return mutes, nil
}

func (d *SqliteDatabase) DeleteMute(mute *types.Mute) (bool, error) {
	query := "DELETE FROM mutes WHERE"
	args := []any{}

	if mute.Chain.IsZero() && mute.ProposalID.IsZero() {
		query += " chain IS NULL AND proposal_id IS NULL"
	} else if mute.Chain.IsZero() {
		query += " chain IS NULL AND proposal_id = $1"
		args = append(args, mute.ProposalID.String)
	} else if mute.ProposalID.IsZero() {
		query += " chain = $1 AND proposal_id IS NULL"
		args = append(args, mute.Chain.String)
	} else {
		query += " chain = $1 AND proposal_id = $2"
		args = append(args, mute.Chain.String, mute.ProposalID.String)
	}

	result, err := d.client.Exec(query, args...)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete mute")
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not get affected rows on deleting mute")
		return false, err
	}

	return rowsAffected > 0, nil
}

func (d *SqliteDatabase) IsMuted(chain, proposalID string) (bool, error) {
	row := d.client.QueryRow(
		"SELECT COUNT(*) FROM mutes WHERE ((chain IS NULL AND proposal_id IS NULL) OR (chain = $1 AND proposal_id IS NULL) OR (chain IS NULL AND proposal_id = $2) OR (chain = $1 AND proposal_id = $2)) AND expires >= datetime('now')",
		chain,
		proposalID,
	)

	if err := row.Err(); err != nil {
		d.logger.Error().Err(err).Msg("Could not check if entry was muted")
		return false, err
	}

	count := 0

	if err := row.Scan(&count); err != nil {
		d.logger.Error().Err(err).Msg("Could not scan to check entry was muted")
		return false, err
	}

	return count > 0, nil
}

func (d *SqliteDatabase) Destroy() error {
	return os.Remove(d.config.Path)
}
