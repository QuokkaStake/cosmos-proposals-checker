package cosmos

import (
	"encoding/json"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"net/http"
	"time"

	"main/pkg/types"

	"github.com/rs/zerolog"
)

const PaginationLimit = 1000

type RPC struct {
	ChainConfig   *types.Chain
	URLs          []string
	ProposalsType string
	Logger        zerolog.Logger
}

func NewRPC(chainConfig *types.Chain, logger zerolog.Logger) *RPC {
	return &RPC{
		ChainConfig:   chainConfig,
		URLs:          chainConfig.LCDEndpoints,
		ProposalsType: chainConfig.ProposalsType,
		Logger:        logger.With().Str("component", "rpc").Logger(),
	}
}

func (rpc *RPC) GetAllProposals() ([]types.Proposal, *types.QueryError) {
	if rpc.ProposalsType == "v1" {
		return rpc.GetAllV1Proposals()
	}

	return rpc.GetAllV1beta1Proposals()
}

func (rpc *RPC) GetTally(proposal string) (*types.Tally, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/tally",
		proposal,
	)

	var tally responses.TallyRPCResponse
	if errs := rpc.Get(url, &tally); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return tally.Tally.ToTally(), nil
}

func (rpc *RPC) GetStakingPool() (*types.PoolRPCResponse, *types.QueryError) {
	url := "/cosmos/staking/v1beta1/pool"

	var pool types.PoolRPCResponse
	if errs := rpc.Get(url, &pool); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}

func (rpc *RPC) GetGovParams(paramsType string) (*types.ParamsResponse, *types.QueryError) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/params/%s", paramsType)

	var pool types.ParamsResponse
	if errs := rpc.Get(url, &pool); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}

func (rpc *RPC) Get(url string, target interface{}) []types.NodeError {
	nodeErrors := make([]types.NodeError, len(rpc.URLs))

	for index, lcd := range rpc.URLs {
		fullURL := lcd + url
		rpc.Logger.Trace().Str("url", fullURL).Msg("Trying making request to LCD")

		err := rpc.GetFull(
			fullURL,
			target,
		)

		if err == nil {
			return nil
		}

		rpc.Logger.Warn().Str("url", fullURL).Err(err).Msg("LCD request failed")
		nodeErrors[index] = types.NodeError{
			Node:  lcd,
			Error: types.NewJSONError(err),
		}
	}

	rpc.Logger.Warn().Str("url", url).Msg("All LCD requests failed")
	return nodeErrors
}

func (rpc *RPC) GetFull(url string, target interface{}) error {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "cosmos-proposals-checker")

	rpc.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		rpc.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	rpc.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}
