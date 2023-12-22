package cosmos

import (
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"sync"
)

func (rpc *RPC) GetGovParams(paramsType string) (*responses.ParamsResponse, *types.QueryError) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/params/%s", paramsType)

	var params responses.ParamsResponse
	if errs := rpc.Client.Get(url, &params); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &params, nil
}

func (rpc *RPC) GetChainParams() (*types.ChainWithVotingParams, []error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)
	params := &responses.ParamsResponse{}

	wg.Add(3)

	go func() {
		defer wg.Done()

		votingParams, err := rpc.GetGovParams("voting")
		mutex.Lock()
		defer mutex.Unlock()

		if err != nil {
			errors = append(errors, err)
			return
		}

		params.VotingParams = votingParams.VotingParams
	}()

	go func() {
		defer wg.Done()

		depositParams, err := rpc.GetGovParams("deposit")
		mutex.Lock()
		defer mutex.Unlock()

		if err != nil {
			errors = append(errors, err)
			return
		}

		params.DepositParams = depositParams.DepositParams
	}()

	go func() {
		defer wg.Done()

		tallyingParams, err := rpc.GetGovParams("tallying")
		mutex.Lock()
		defer mutex.Unlock()

		if err != nil {
			errors = append(errors, err)
			return
		}

		params.TallyParams = tallyingParams.TallyParams
	}()

	wg.Wait()

	if len(errors) > 0 {
		return nil, errors
	}

	return params.ToParams(rpc.ChainConfig)
}
