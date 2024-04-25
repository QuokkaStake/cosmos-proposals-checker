package neutron

import (
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetAllProposals() ([]types.Proposal, *types.QueryError) {
	query := "{\"list_proposals\": {}}"

	var proposals responses.ProposalsResponse
	if _, err := fetcher.GetSmartContractState(query, &proposals, 0); err != nil {
		return nil, err
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
