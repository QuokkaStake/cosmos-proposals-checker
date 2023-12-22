package cosmos

import (
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/types"
	"main/pkg/utils"
)

func (rpc *RPC) GetAllV1beta1Proposals() ([]types.Proposal, *types.QueryError) {
	proposals := []types.Proposal{}
	offset := 0

	for {
		url := fmt.Sprintf(
			// 2 is for PROPOSAL_STATUS_VOTING_PERIOD
			"/cosmos/gov/v1beta1/proposals?pagination.limit=%d&pagination.offset=%d&proposal_status=2",
			PaginationLimit,
			offset,
		)

		var batchProposals responses.V1Beta1ProposalsRPCResponse
		if errs := rpc.Client.Get(url, &batchProposals); len(errs) > 0 {
			return nil, &types.QueryError{
				QueryError: nil,
				NodeErrors: errs,
			}
		}

		if batchProposals.Message != "" {
			return nil, &types.QueryError{
				QueryError: errors.New(batchProposals.Message),
			}
		}

		parsedProposals := utils.Map(batchProposals.Proposals, func(p responses.V1beta1Proposal) types.Proposal {
			return p.ToProposal()
		})
		proposals = append(proposals, parsedProposals...)
		if len(batchProposals.Proposals) < PaginationLimit {
			break
		}

		offset += PaginationLimit
	}

	return proposals, nil
}
