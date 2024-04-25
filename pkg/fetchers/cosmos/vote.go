package cosmos

import (
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
)

func (rpc *RPC) GetVote(proposal, voter string, prevHeight int64) (*types.Vote, int64, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/votes/%s",
		proposal,
		voter,
	)

	var vote responses.VoteRPCResponse
	errs, header := rpc.Client.GetWithPredicate(
		url,
		&vote,
		types.HTTPPredicateCheckHeightAfter(prevHeight),
	)
	if len(errs) > 0 {
		return nil, 0, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	height, err := utils.GetBlockHeightFromHeader(header)
	if err != nil {
		return nil, 0, &types.QueryError{
			QueryError: errors.New("got error when parsing vote height"),
		}
	}

	if vote.IsError() {
		// not voted
		if strings.Contains(vote.Message, "not found") {
			return nil, height, nil
		}

		// some other errors
		return nil, height, &types.QueryError{
			QueryError: errors.New(vote.Message),
		}
	}

	voteParsed, err := vote.ToVote()
	if err != nil {
		return nil, height, &types.QueryError{
			QueryError: err,
			NodeErrors: nil,
		}
	}

	return voteParsed, height, nil
}
