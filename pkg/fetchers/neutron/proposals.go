package neutron

import (
	"encoding/base64"
	"fmt"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetAllProposals() ([]types.Proposal, *types.QueryError) {
	query := base64.StdEncoding.EncodeToString([]byte("{\"list_proposals\": {}}"))

	url := fmt.Sprintf(
		"/cosmwasm/wasm/v1/contract/%s/smart/%s",
		fetcher.ChainConfig.NeutronSmartContract,
		query,
	)

	var proposals responses.ProposalsResponse
	if errs := fetcher.Client.Get(url, &proposals); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	proposalsParsed, err := proposals.ToProposals()
	if err != nil {
		return nil, &types.QueryError{
			QueryError: err,
			NodeErrors: nil,
		}
	}

	return proposalsParsed, nil
}
