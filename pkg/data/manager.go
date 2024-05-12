package data

import (
	"fmt"
	fetchersPkg "main/pkg/fetchers"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Manager struct {
	Logger   zerolog.Logger
	Chains   types.Chains
	Fetchers []fetchersPkg.Fetcher
}

func NewManager(logger *zerolog.Logger, chains types.Chains) *Manager {
	fetchers := make([]fetchersPkg.Fetcher, len(chains))

	for index, chain := range chains {
		fetchers[index] = fetchersPkg.GetFetcher(chain, logger)
	}

	return &Manager{
		Logger:   logger.With().Str("component", "data_manager").Logger(),
		Chains:   chains,
		Fetchers: fetchers,
	}
}

func (m *Manager) GetTallies() (map[string]types.ChainTallyInfos, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)
	tallies := make(map[string]types.ChainTallyInfos)

	for index, chain := range m.Chains {
		fetcher := m.Fetchers[index]

		wg.Add(1)
		go func(c *types.Chain, fetcher fetchersPkg.Fetcher) {
			defer wg.Done()

			talliesForChain, err := fetcher.GetTallies()

			mutex.Lock()

			if err != nil {
				m.Logger.Error().Err(err).Str("chain", c.Name).Msg("Error fetching tallies")
				errors = append(errors, err)
			} else if len(talliesForChain.TallyInfos) > 0 {
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

	for index := range m.Chains {
		fetcher := m.Fetchers[index]

		wg.Add(1)

		go func(fetcher fetchersPkg.Fetcher) {
			defer wg.Done()

			chainParams, errs := fetcher.GetChainParams()
			mutex.Lock()
			defer mutex.Unlock()

			if len(errs) > 0 {
				errors = append(errors, errs...)
				return
			}

			params[chainParams.Chain.Name] = *chainParams
		}(fetcher)
	}

	wg.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("got %d errors when fetching chain params", len(errors))
	}

	return params, nil
}
