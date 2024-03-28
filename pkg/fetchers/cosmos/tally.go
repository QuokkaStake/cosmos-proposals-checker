package cosmos

import (
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"sync"

	"cosmossdk.io/math"
)

func (rpc *RPC) GetTally(proposal string) (*types.Tally, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/tally",
		proposal,
	)

	var tally responses.TallyRPCResponse
	if errs := rpc.Client.Get(url, &tally); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return tally.Tally.ToTally(), nil
}

func (rpc *RPC) GetTallies() (types.ChainTallyInfos, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	errorsList := make([]error, 0)

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
			errorsList = append(errorsList, err)
		} else if poolResponse.Pool == nil {
			rpc.Logger.Error().Err(err).Msg("Staking pool is empty!")
			errorsList = append(errorsList, errors.New("staking pool is empty"))
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
			errorsList = append(errorsList, err)

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
					errorsList = append(errorsList, err)
				} else if tally == nil {
					rpc.Logger.Error().
						Err(err).
						Str("proposal_id", p.ID).
						Msg("Tally is empty")
					errorsList = append(errorsList, errors.New("tally is empty"))
				} else {
					tallies[p.ID] = *tally
				}
			}(proposal)
		}

		internalWg.Wait()
	}()

	wg.Wait()

	if len(errorsList) > 0 {
		rpc.Logger.Error().Msg("Errors getting tallies info, not processing")
		return types.ChainTallyInfos{}, fmt.Errorf("could not get tallies info: got %d errors", len(errorsList))
	}

	tallyInfos := types.ChainTallyInfos{
		Chain:      rpc.ChainConfig,
		TallyInfos: make([]types.TallyInfo, len(proposals)),
	}

	for index, proposal := range proposals {
		tally, ok := tallies[proposal.ID]
		if !ok {
			return types.ChainTallyInfos{}, errors.New("could not get tallies info")
		}

		tallyInfos.TallyInfos[index] = types.TallyInfo{
			Proposal:         proposal,
			Tally:            tally,
			TotalVotingPower: pool,
		}
	}

	return tallyInfos, nil
}
