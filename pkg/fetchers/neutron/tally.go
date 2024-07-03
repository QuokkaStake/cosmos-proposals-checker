package neutron

import (
	"context"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetTallies(ctx context.Context) (types.ChainTallyInfos, error) {
	query := "{\"reverse_proposals\": {\"limit\": 1000}}"

	var proposals responses.ProposalsResponse
	if _, err := fetcher.GetSmartContractState(query, &proposals, 0, ctx); err != nil {
		return types.ChainTallyInfos{}, err
	}

	return types.ChainTallyInfos{
		Chain:      fetcher.ChainConfig,
		TallyInfos: proposals.ToTally(),
	}, nil
}
