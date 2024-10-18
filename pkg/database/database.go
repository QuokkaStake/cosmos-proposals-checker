package database

import (
	"context"
	"main/pkg/types"
)

type Database interface {
	Init()
	Migrate()
	Rollback()
	UpsertProposal(chain *types.Chain, proposal types.Proposal) error
	GetProposal(chain *types.Chain, proposalID string) (*types.Proposal, error)
	GetVote(chain *types.Chain, proposal types.Proposal, wallet *types.Wallet) (*types.Vote, error)
	UpsertVote(
		chain *types.Chain,
		proposal types.Proposal,
		wallet *types.Wallet,
		vote *types.Vote,
		ctx context.Context,
	) error
	GetLastBlockHeight(chain *types.Chain, storableKey string) (int64, error)
	UpsertLastBlockHeight(chain *types.Chain, storableKey string, height int64) error
	UpsertMute(mute *types.Mute) error
	DeleteMute(mute *types.Mute) (bool, error)
	GetAllMutes() ([]*types.Mute, error)
}
