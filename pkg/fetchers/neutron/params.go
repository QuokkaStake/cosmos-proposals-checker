package neutron

import (
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetChainParams() (*types.ChainWithVotingParams, []error) {
	query := "{\"config\":{}}"

	var params responses.ParamsResponse
	if err := fetcher.GetSmartContractState(query, &params); err != nil {
		return nil, []error{err}
	}

	paramsParsed, errs := params.ToParams(fetcher.ChainConfig)
	if len(errs) > 0 {
		return nil, errs
	}

	return paramsParsed, nil
}
