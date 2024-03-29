package neutron

import (
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetTallies() (types.ChainTallyInfos, error) {
	query := "{\"list_proposals\": {}}"

	var proposals responses.ProposalsResponse
	if err := fetcher.GetSmartContractState(query, &proposals); err != nil {
		return types.ChainTallyInfos{}, err
	}

	tallyInfos, err := proposals.ToTally()
	if err != nil {
		return types.ChainTallyInfos{}, err
	}

	return types.ChainTallyInfos{
		Chain:      fetcher.ChainConfig,
		TallyInfos: tallyInfos,
	}, nil
}
