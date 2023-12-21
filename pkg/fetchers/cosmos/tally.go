package cosmos

import (
	"fmt"
	"main/pkg/types"
	"sync"

	"cosmossdk.io/math"
)

func (rpc *RPC) GetTallies() (types.ChainTallyInfos, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errors := make([]error, 0)

	var pool math.LegacyDec
	var proposals []types.Proposal
	tallies := make(map[string]types.Tally)

	wg.Add(1)
	go func() {
		defer wg.Done()

		poolResponse, err := rpc.GetStakingPool()

		mutex.Lock()

		if err != nil {
			rpc.Logger.Error().Err(err).Msg("Error fetching staking pool")
			errors = append(errors, err)
		} else if poolResponse.Pool == nil {
			rpc.Logger.Error().Err(err).Msg("Staking pool is empty!")
			errors = append(errors, fmt.Errorf("staking pool is empty"))
		} else {
			pool = poolResponse.Pool.BondedTokens
		}
		mutex.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		chainProposals, err := rpc.GetAllProposals()

		mutex.Lock()

		if err != nil {
			rpc.Logger.Error().Err(err).Msg("Error fetching chain proposals")
			errors = append(errors, err)

			mutex.Unlock()
			return
		} else {
			proposals = chainProposals
		}

		mutex.Unlock()

		var internalWg sync.WaitGroup

		for _, proposal := range chainProposals {
			internalWg.Add(1)

			go func(p types.Proposal) {
				defer internalWg.Done()

				tally, err := rpc.GetTally(p.ID)

				mutex.Lock()
				defer mutex.Unlock()

				if err != nil {
					rpc.Logger.Error().
						Err(err).
						Str("proposal_id", p.ID).
						Msg("Error fetching tally for proposal")
					errors = append(errors, err)
				} else if tally == nil {
					rpc.Logger.Error().
						Err(err).
						Str("proposal_id", p.ID).
						Msg("Tally is empty")
					errors = append(errors, fmt.Errorf("tally is empty"))
				} else {
					tallies[p.ID] = *tally
				}
			}(proposal)
		}

		internalWg.Wait()
	}()

	wg.Wait()

	if len(errors) > 0 {
		rpc.Logger.Error().Msg("Errors getting tallies info, not processing")
		return types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info: got %d errors", len(errors))
	}

	tallyInfos := types.ChainTallyInfos{
		Chain:      rpc.ChainConfig,
		TallyInfos: make([]types.TallyInfo, len(proposals)),
	}

	for index, proposal := range proposals {
		tally, ok := tallies[proposal.ID]
		if !ok {
			return types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info")
		}

		tallyInfos.TallyInfos[index] = types.TallyInfo{
			Proposal:         proposal,
			Tally:            tally,
			TotalVotingPower: pool,
		}
	}

	return tallyInfos, nil
}
