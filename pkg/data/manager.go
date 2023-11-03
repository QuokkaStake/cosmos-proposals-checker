package data

import (
	"fmt"
	"github.com/rs/zerolog"
	configTypes "main/pkg/config/types"
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"
)

type Manager struct {
	Logger zerolog.Logger
	Chains configTypes.Chains
}

func NewManager(logger *zerolog.Logger, chains configTypes.Chains) *Manager {
	return &Manager{
		Logger: logger.With().Str("component", "data_manager").Logger(),
		Chains: chains,
	}
}

func (m *Manager) GetTallies() (map[string][]types.TallyInfo, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)

	pools := make(map[string]types.Pool, 0)
	proposals := make(map[string][]types.Proposal, 0)
	tallies := make(map[string]map[string]types.Tally, 0)

	for _, chain := range m.Chains {
		rpc := tendermint.NewRPC(chain, m.Logger)

		wg.Add(1)
		go func(c *configTypes.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			pool, err := rpc.GetStakingPool()

			mutex.Lock()

			if err != nil {
				m.Logger.Error().Err(err).Str("chain", c.Name).Msg("Error fetching staking pool")
				errors = append(errors, err)
			} else if pool.Pool == nil {
				m.Logger.Error().Err(err).Str("chain", c.Name).Msg("Staking pool is empty!")
				errors = append(errors, fmt.Errorf("staking pool is empty"))
			} else {
				pools[c.Name] = *pool.Pool
			}
			mutex.Unlock()

		}(chain, rpc)

		wg.Add(1)
		go func(c *configTypes.Chain, rpc *tendermint.RPC) {
			defer wg.Done()

			chainProposals, err := rpc.GetAllProposals()

			mutex.Lock()

			if err != nil {
				m.Logger.Error().Err(err).Str("chain", c.Name).Msg("Error fetching chain proposals")
				errors = append(errors, err)

				mutex.Unlock()
				return
			} else {
				proposals[c.Name] = chainProposals
			}

			mutex.Unlock()

			var internalWg sync.WaitGroup

			for _, proposal := range chainProposals {
				internalWg.Add(1)

				go func(c *configTypes.Chain, p types.Proposal) {
					defer internalWg.Done()

					tally, err := rpc.GetTally(p.ID)

					mutex.Lock()
					defer mutex.Unlock()

					if err != nil {
						m.Logger.Error().
							Err(err).
							Str("chain", c.Name).
							Str("proposal_id", p.ID).
							Msg("Error fetching tally for proposal")
						errors = append(errors, err)
					} else if tally.Tally == nil {
						m.Logger.Error().
							Err(err).
							Str("chain", c.Name).
							Str("proposal_id", p.ID).
							Msg("Tally is empty")
						errors = append(errors, fmt.Errorf("tally is empty"))
					} else {
						if _, ok := tallies[c.Name]; !ok {
							tallies[c.Name] = make(map[string]types.Tally, 0)
						}

						tallies[c.Name][p.ID] = *tally.Tally
					}
				}(c, proposal)
			}

			internalWg.Wait()
		}(chain, rpc)
	}

	wg.Wait()

	if len(errors) > 0 {
		m.Logger.Error().Msg("Errors getting tallies info, not processing")
		return map[string][]types.TallyInfo{}, fmt.Errorf("could not get tallies info")
	}

	tallyInfos := make(map[string][]types.TallyInfo, 0)

	for chainName, chainProposals := range proposals {
		for _, proposal := range chainProposals {
			tally, ok := tallies[chainName][proposal.ID]
			if !ok {
				return map[string][]types.TallyInfo{}, fmt.Errorf("could not get tallies info")
			}

			pool, ok := pools[chainName]
			if !ok {
				return map[string][]types.TallyInfo{}, fmt.Errorf("could not get tallies info")
			}

			if _, ok := tallyInfos[chainName]; !ok {
				tallyInfos[chainName] = []types.TallyInfo{}
			}

			tallyInfos[chainName] = append(tallyInfos[chainName], types.TallyInfo{
				ChainName: chainName,
				Proposal:  proposal,
				Tally:     tally,
				Pool:      pool,
			})
		}
	}

	return tallyInfos, nil
}
