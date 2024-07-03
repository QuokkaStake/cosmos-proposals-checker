package cosmos

import (
	"context"
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"main/pkg/utils"
)

func (rpc *RPC) GetAllV1beta1Proposals(
	prevHeight int64,
	ctx context.Context,
) ([]types.Proposal, int64, *types.QueryError) {
	proposals := []types.Proposal{}
	offset := 0

	lastHeight := prevHeight

	for {
		url := fmt.Sprintf(
			"/cosmos/gov/v1beta1/proposals?pagination.limit=%d&pagination.offset=%d&pagination.count_total=1",
			rpc.PaginationLimit,
			offset,
		)

		var batchProposals responses.V1Beta1ProposalsRPCResponse
		errs, header := rpc.Client.GetWithPredicate(
			url,
			&batchProposals,
			types.HTTPPredicateCheckHeightAfter(lastHeight),
			ctx,
		)
		if len(errs) > 0 {
			return nil, 0, &types.QueryError{
				QueryError: nil,
				NodeErrors: errs,
			}
		}

		height, _ := utils.GetBlockHeightFromHeader(header)

		if batchProposals.Message != "" {
			return nil, height, &types.QueryError{
				QueryError: errors.New(batchProposals.Message),
			}
		}

		lastHeight = height

		parsedProposals := utils.Map(batchProposals.Proposals, func(p responses.V1beta1Proposal) types.Proposal {
			return p.ToProposal()
		})
		proposals = append(proposals, parsedProposals...)
		if len(batchProposals.Proposals) < rpc.PaginationLimit {
			break
		}

		offset += rpc.PaginationLimit
	}

	return proposals, lastHeight, nil
}
