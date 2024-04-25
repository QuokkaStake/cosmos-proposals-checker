package neutron

import (
	"fmt"
	"main/pkg/fetchers/neutron/responses"
	"main/pkg/types"
)

func (fetcher *Fetcher) GetVote(
	proposal, voter string,
	prevHeight int64,
) (*types.Vote, int64, *types.QueryError) {
	query := fmt.Sprintf(
		"{\"get_vote\":{\"proposal_id\":%s,\"voter\":\"%s\"}}",
		proposal,
		voter,
	)

	var vote responses.VoteResponse
	height, err := fetcher.GetSmartContractState(query, &vote, prevHeight)
	if err != nil {
		return nil, 0, err
	}

	voteParsed := vote.ToVote(proposal)
	return voteParsed, height, nil
}
