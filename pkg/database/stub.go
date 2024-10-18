package database

import (
	"context"
	"main/pkg/types"
)

type StubDatabase struct {
	LastHeightQueryErrors map[string]map[string]error
	LastHeightWriteError  error
	GetProposalError      error
	UpsertProposalError   error
	GetVoteError          error
	UpsertVoteError       error

	Proposals       map[string]map[string]*types.Proposal
	Votes           map[string]map[string]map[string]*types.Vote
	LastBlockHeight map[string]map[string]int64
}

func (d *StubDatabase) Init() {

}

func (d *StubDatabase) Migrate() {

}

func (d *StubDatabase) Rollback() {

}

func (d *StubDatabase) UpsertProposal(
	chain *types.Chain,
	proposal types.Proposal,
) error {
	if d.UpsertProposalError != nil {
		return d.UpsertProposalError
	}

	if d.Proposals == nil {
		d.Proposals = make(map[string]map[string]*types.Proposal)
	}

	if _, ok := d.Proposals[chain.Name]; !ok {
		d.Proposals[chain.Name] = make(map[string]*types.Proposal)
	}

	d.Proposals[chain.Name][proposal.ID] = &proposal
	return nil
}

func (d *StubDatabase) GetProposal(chain *types.Chain, proposalID string) (*types.Proposal, error) {
	if d.GetProposalError != nil {
		return nil, d.GetProposalError
	}

	chainProposals, ok := d.Proposals[chain.Name]
	if !ok {
		return nil, nil //nolint:nilnil
	}

	return chainProposals[proposalID], nil
}

func (d *StubDatabase) GetVote(
	chain *types.Chain,
	proposal types.Proposal,
	wallet *types.Wallet,
) (*types.Vote, error) {
	if d.GetVoteError != nil {
		return nil, d.GetVoteError
	}

	chainVotes, ok := d.Votes[chain.Name]
	if !ok {
		return nil, nil //nolint:nilnil
	}

	proposalVotes, ok := chainVotes[proposal.ID]
	if !ok {
		return nil, nil //nolint:nilnil
	}

	return proposalVotes[wallet.Address], nil
}

func (d *StubDatabase) UpsertVote(
	chain *types.Chain,
	proposal types.Proposal,
	wallet *types.Wallet,
	vote *types.Vote,
	ctx context.Context,
) error {
	if d.UpsertVoteError != nil {
		return d.UpsertVoteError
	}

	if _, ok := d.Votes[chain.Name]; !ok {
		d.Votes[chain.Name] = make(map[string]map[string]*types.Vote)
	}

	if _, ok := d.Votes[chain.Name][proposal.ID]; !ok {
		d.Votes[chain.Name][proposal.ID] = make(map[string]*types.Vote)
	}

	d.Votes[chain.Name][proposal.ID][wallet.Address] = vote
	return nil
}

func (d *StubDatabase) GetLastBlockHeight(
	chain *types.Chain,
	storableKey string,
) (int64, error) {
	if d.LastHeightQueryErrors != nil {
		if chainErrors, chainErrorsFound := d.LastHeightQueryErrors[chain.Name]; chainErrorsFound {
			if err, errFound := chainErrors[storableKey]; errFound {
				return 0, err
			}
		}
	}

	chainHeights, ok := d.LastBlockHeight[chain.Name]
	if !ok {
		return 0, nil
	}

	return chainHeights[storableKey], nil
}

func (d *StubDatabase) UpsertLastBlockHeight(
	chain *types.Chain,
	storableKey string,
	height int64,
) error {
	if d.LastHeightWriteError != nil {
		return d.LastHeightWriteError
	}

	if d.LastBlockHeight == nil {
		d.LastBlockHeight = make(map[string]map[string]int64)
	}

	if _, ok := d.LastBlockHeight[chain.Name]; !ok {
		d.LastBlockHeight[chain.Name] = make(map[string]int64)
	}

	d.LastBlockHeight[chain.Name][storableKey] = height
	return nil
}

func (d *StubDatabase) UpsertMute(mute *types.Mute) error {
	return nil
}

func (d *StubDatabase) GetAllMutes() ([]*types.Mute, error) {
	return []*types.Mute{}, nil
}

func (d *StubDatabase) DeleteMute(mute *types.Mute) (bool, error) {
	return true, nil
}
