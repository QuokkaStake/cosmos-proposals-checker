package data

import (
	"fmt"
	"main/pkg/fetchers"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Manager struct {
	Logger zerolog.Logger
	Chains types.Chains
}

func NewManager(logger *zerolog.Logger, chains types.Chains) *Manager {
	return &Manager{
		Logger: logger.With().Str("component", "data_manager").Logger(),
		Chains: chains,
	}
}

func (m *Manager) GetTallies() (map[string]types.ChainTallyInfos, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)
	tallies := make(map[string]types.ChainTallyInfos)

	for _, chain := range m.Chains {
		fetcher := fetchers.GetFetcher(chain, m.Logger)

		wg.Add(1)
		go func(c *types.Chain, fetcher fetchers.Fetcher) {
			defer wg.Done()

			talliesForChain, err := fetcher.GetTallies()

			mutex.Lock()

			if err != nil {
				m.Logger.Error().Err(err).Str("chain", c.Name).Msg("Error fetching staking pool")
				errors = append(errors, err)
			} else {
				tallies[c.Name] = talliesForChain
			}
			mutex.Unlock()
		}(chain, fetcher)
	}

	wg.Wait()

	if len(errors) > 0 {
		m.Logger.Error().Msg("Errors getting tallies info, not processing")
		return map[string]types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info: got %d errors", len(errors))
	}

	return tallies, nil
}

func (m *Manager) GetParams() (map[string]types.ChainWithVotingParams, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	params := make(map[string]types.ChainWithVotingParams)
	errors := make([]error, 0)

	for _, chain := range m.Chains {
		wg.Add(1)

		go func(chain *types.Chain) {
			defer wg.Done()

			fetcher := fetchers.GetFetcher(chain, m.Logger)

			chainParams, errs := fetcher.GetChainParams()
			mutex.Lock()
			defer mutex.Unlock()

			if len(errs) > 0 {
				errors = append(errors, errs...)
				return
			}

			params[chainParams.Chain.Name] = *chainParams
		}(chain)
	}

	wg.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("got %d errors when fetching chain params", len(errors))
	}

	return params, nil
}
