package data

import (
	"fmt"
	"main/pkg/tendermint"
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

	pools := make(map[string]types.Pool, 0)
	proposals := make(map[string][]types.Proposal, 0)
	tallies := make(map[string]map[string]types.Tally, 0)

	for _, chain := range m.Chains {
		rpc := tendermint.NewRPC(chain, m.Logger)

		wg.Add(1)
		go func(c *types.Chain, rpc *tendermint.RPC) {
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
		go func(c *types.Chain, rpc *tendermint.RPC) {
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

				go func(c *types.Chain, p types.Proposal) {
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
		return map[string]types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info")
	}

	tallyInfos := make(map[string]types.ChainTallyInfos, 0)

	for chainName, chainProposals := range proposals {
		chain := m.Chains.FindByName(chainName)
		if chain == nil {
			return map[string]types.ChainTallyInfos{}, fmt.Errorf("could not chain with name %s", chainName)
		}

		if _, ok := tallyInfos[chainName]; !ok {
			tallyInfos[chainName] = types.ChainTallyInfos{
				Chain:      chain,
				TallyInfos: make([]types.TallyInfo, len(chainProposals)),
			}
		}

		for index, proposal := range chainProposals {
			tally, ok := tallies[chainName][proposal.ID]
			if !ok {
				return map[string]types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info")
			}

			pool, ok := pools[chainName]
			if !ok {
				return map[string]types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info")
			}

			tallyInfos[chainName].TallyInfos[index] = types.TallyInfo{
				Proposal: proposal,
				Tally:    tally,
				Pool:     pool,
			}
		}
	}

	return tallyInfos, nil
}
