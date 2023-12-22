package neutron

import "main/pkg/types"

func (fetcher *Fetcher) GetChainParams() (*types.ChainWithVotingParams, []error) {
	// TODO: fix
	return &types.ChainWithVotingParams{
		Chain:  fetcher.ChainConfig,
		Params: []types.ChainParam{},
	}, nil
}
