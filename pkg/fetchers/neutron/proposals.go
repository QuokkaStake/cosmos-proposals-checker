package neutron

import (
	"context"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetAllProposals(
	prevHeight int64,
	ctx context.Context,
) ([]types.Proposal, int64, *types.QueryError) {
	query := "{\"reverse_proposals\": {\"limit\": 1000}}"

	var proposals responses.ProposalsResponse
	height, err := fetcher.GetSmartContractState(query, &proposals, prevHeight, ctx)
	if err != nil {
		return nil, height, err
	}

	return proposals.ToProposals(), height, nil
}
