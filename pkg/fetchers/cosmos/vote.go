package cosmos

import (
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"strings"
)

func (rpc *RPC) GetVote(proposal, voter string) (*types.Vote, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/votes/%s",
		proposal,
		voter,
	)

	var vote responses.VoteRPCResponse
	if errs := rpc.Client.Get(url, &vote); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	if vote.IsError() {
		// not voted
		if strings.Contains(vote.Message, "not found") {
			return nil, nil
		}

		// some other errors
		return nil, &types.QueryError{
			QueryError: errors.New(vote.Message),
		}
	}

	voteParsed, err := vote.ToVote()
	if err != nil {
		return nil, &types.QueryError{
			QueryError: err,
			NodeErrors: nil,
		}
	}

	return voteParsed, nil
}
