package neutron

import (
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetAllProposals(
	prevHeight int64,
) ([]types.Proposal, int64, *types.QueryError) {
	query := "{\"list_proposals\": {}}"

	var proposals responses.ProposalsResponse
	height, err := fetcher.GetSmartContractState(query, &proposals, prevHeight)
	if err != nil {
		return nil, height, err
	}

	proposalsParsed, parseErr := proposals.ToProposals()
	if parseErr != nil {
		return nil, height, &types.QueryError{
			QueryError: err,
			NodeErrors: nil,
		}
	}

	return proposalsParsed, height, nil
}
