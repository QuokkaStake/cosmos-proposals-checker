package neutron

import (
	"context"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetChainParams(ctx context.Context) (*types.ChainWithVotingParams, []error) {
	query := "{\"config\":{}}"

	var params responses.ParamsResponse
	if _, err := fetcher.GetSmartContractState(query, &params, 0, ctx); err != nil {
		return nil, []error{err}
	}

	paramsParsed, errs := params.ToParams(fetcher.ChainConfig)
	if len(errs) > 0 {
		return nil, errs
	}

	return paramsParsed, nil
}
