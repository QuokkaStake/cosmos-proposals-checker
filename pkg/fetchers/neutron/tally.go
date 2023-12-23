package neutron

import (
	"encoding/base64"
	"fmt"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetTallies() (types.ChainTallyInfos, error) {
	query := base64.StdEncoding.EncodeToString([]byte("{\"list_proposals\": {}}"))

	url := fmt.Sprintf(
		"/cosmwasm/wasm/v1/contract/%s/smart/%s",
		fetcher.ChainConfig.NeutronSmartContract,
		query,
	)

	var proposals responses.ProposalsResponse
	if errs := fetcher.Client.Get(url, &proposals); len(errs) > 0 {
		return types.ChainTallyInfos{}, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
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
