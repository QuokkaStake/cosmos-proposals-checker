package cosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/utils"
	"net/http"
	"strings"
	"time"

	"main/pkg/types"

	"github.com/rs/zerolog"
)

const PaginationLimit = 1000

type RPC struct {
	URLs          []string
	ProposalsType string
	Logger        zerolog.Logger
}

func NewRPC(chainConfig *types.Chain, logger zerolog.Logger) *RPC {
	return &RPC{
		URLs:          chainConfig.LCDEndpoints,
		ProposalsType: chainConfig.ProposalsType,
		Logger:        logger.With().Str("component", "rpc").Logger(),
	}
}

func (rpc *RPC) GetAllProposals() ([]types.Proposal, *types.QueryError) {
	if rpc.ProposalsType == "v1" {
		return rpc.GetAllV1Proposals()
	}

	return rpc.GetAllV1beta1Proposals()
}

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

		var batchProposals types.V1Beta1ProposalsRPCResponse
		if errs := rpc.Get(url, &batchProposals); len(errs) > 0 {
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

		parsedProposals := utils.Map(batchProposals.Proposals, func(p types.V1beta1Proposal) types.Proposal {
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

func (rpc *RPC) GetAllV1Proposals() ([]types.Proposal, *types.QueryError) {
	proposals := []types.Proposal{}
	offset := 0

	for {
		url := fmt.Sprintf(
			// 2 is for PROPOSAL_STATUS_VOTING_PERIOD
			"/cosmos/gov/v1/proposals?pagination.limit=%d&pagination.offset=%d&proposal_status=2",
			PaginationLimit,
			offset,
		)

		var batchProposals types.V1ProposalsRPCResponse
		if errs := rpc.Get(url, &batchProposals); len(errs) > 0 {
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

		parsedProposals := utils.Map(batchProposals.Proposals, func(p types.V1Proposal) types.Proposal {
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

func (rpc *RPC) GetVote(proposal, voter string) (*types.Vote, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/votes/%s",
		proposal,
		voter,
	)

	var vote responses.VoteRPCResponse
	if errs := rpc.Get(url, &vote); len(errs) > 0 {
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

func (rpc *RPC) GetTally(proposal string) (*types.TallyRPCResponse, *types.QueryError) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/tally",
		proposal,
	)

	var tally types.TallyRPCResponse
	if errs := rpc.Get(url, &tally); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &tally, nil
}

func (rpc *RPC) GetStakingPool() (*types.PoolRPCResponse, *types.QueryError) {
	url := "/cosmos/staking/v1beta1/pool"

	var pool types.PoolRPCResponse
	if errs := rpc.Get(url, &pool); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}

func (rpc *RPC) GetGovParams(paramsType string) (*types.ParamsResponse, *types.QueryError) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/params/%s", paramsType)

	var pool types.ParamsResponse
	if errs := rpc.Get(url, &pool); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}

func (rpc *RPC) Get(url string, target interface{}) []types.NodeError {
	nodeErrors := make([]types.NodeError, len(rpc.URLs))

	for index, lcd := range rpc.URLs {
		fullURL := lcd + url
		rpc.Logger.Trace().Str("url", fullURL).Msg("Trying making request to LCD")

		err := rpc.GetFull(
			fullURL,
			target,
		)

		if err == nil {
			return nil
		}

		rpc.Logger.Warn().Str("url", fullURL).Err(err).Msg("LCD request failed")
		nodeErrors[index] = types.NodeError{
			Node:  lcd,
			Error: types.NewJSONError(err),
		}
	}

	rpc.Logger.Warn().Str("url", url).Msg("All LCD requests failed")
	return nodeErrors
}

func (rpc *RPC) GetFull(url string, target interface{}) error {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "cosmos-proposals-checker")

	rpc.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		rpc.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	rpc.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}
