package neutron

import (
	"encoding/base64"
	"errors"
	"fmt"
	"main/pkg/http"
	"main/pkg/types"
	"main/pkg/utils"

	"github.com/rs/zerolog"
)

type Fetcher struct {
	ChainConfig *types.Chain
	Logger      zerolog.Logger
	Client      *http.Client
}

func NewFetcher(chainConfig *types.Chain, logger zerolog.Logger) *Fetcher {
	return &Fetcher{
		ChainConfig: chainConfig,
		Logger:      logger.With().Str("component", "neutron_fetcher").Logger(),
		Client:      http.NewClient(chainConfig.Name, chainConfig.LCDEndpoints, logger),
	}
}

func (fetcher *Fetcher) GetSmartContractState(
	queryString string,
	output interface{},
	prevHeight int64,
) (int64, *types.QueryError) {
	query := base64.StdEncoding.EncodeToString([]byte(queryString))

	url := fmt.Sprintf(
		"/cosmwasm/wasm/v1/contract/%s/smart/%s",
		fetcher.ChainConfig.NeutronSmartContract,
		query,
	)

	errs, header := fetcher.Client.GetWithPredicate(
		url,
		&output,
		types.HTTPPredicateCheckHeightAfter(prevHeight),
	)
	if len(errs) > 0 {
		return 0, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	height, err := utils.GetBlockHeightFromHeader(header)
	if err != nil {
		return 0, &types.QueryError{
			QueryError: errors.New("got error when parsing vote height"),
		}
	}

	return height, nil
}
