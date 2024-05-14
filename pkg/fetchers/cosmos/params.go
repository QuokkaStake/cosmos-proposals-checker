package cosmos

import (
	"context"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"sync"
)

func (rpc *RPC) GetGovParams(paramsType string, ctx context.Context) (*responses.ParamsResponse, *types.QueryError) {
	url := "/cosmos/gov/v1beta1/params/" + paramsType

	var params responses.ParamsResponse
	if errs := rpc.Client.Get(url, &params, ctx); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &params, nil
}

func (rpc *RPC) GetChainParams(ctx context.Context) (*types.ChainWithVotingParams, []error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)
	params := &responses.ParamsResponse{}

	wg.Add(3)

	go func() {
		defer wg.Done()

		votingParams, err := rpc.GetGovParams("voting", ctx)
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

		depositParams, err := rpc.GetGovParams("deposit", ctx)
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

		tallyingParams, err := rpc.GetGovParams("tallying", ctx)
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
