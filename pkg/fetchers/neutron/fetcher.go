package neutron

import (
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
