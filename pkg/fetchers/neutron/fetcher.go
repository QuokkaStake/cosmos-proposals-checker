package neutron

import (
	"encoding/base64"
	"fmt"
	"main/pkg/http"
	"main/pkg/types"

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

func (fetcher *Fetcher) GetSmartContractState(queryString string, output interface{}) *types.QueryError {
	query := base64.StdEncoding.EncodeToString([]byte(queryString))

	url := fmt.Sprintf(
		"/cosmwasm/wasm/v1/contract/%s/smart/%s",
		fetcher.ChainConfig.NeutronSmartContract,
		query,
	)

	if errs := fetcher.Client.Get(url, &output); len(errs) > 0 {
		return &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return nil
}
