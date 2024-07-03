package cosmos

import (
	"context"
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
)

func (rpc *RPC) GetVote(
	proposal, voter string,
	prevHeight int64,
	ctx context.Context,
) (*types.Vote, int64, *types.QueryError) {
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
		ctx,
	)
	if len(errs) > 0 {
		return nil, 0, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	height, _ := utils.GetBlockHeightFromHeader(header)

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

	return vote.ToVote(), height, nil
}
