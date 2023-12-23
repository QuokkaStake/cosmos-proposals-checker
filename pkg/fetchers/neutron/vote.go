package neutron

import (
	"fmt"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetVote(proposal, voter string) (*types.Vote, *types.QueryError) {
	query := fmt.Sprintf(
		"{\"get_vote\":{\"proposal_id\":%s,\"voter\":\"%s\"}}",
		proposal,
		voter,
	)

	var vote responses.VoteResponse
	if err := fetcher.GetSmartContractState(query, &vote); err != nil {
		return nil, err
	}

	voteParsed := vote.ToVote(proposal)
	return voteParsed, nil
}
