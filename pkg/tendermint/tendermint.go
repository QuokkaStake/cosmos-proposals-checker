package tendermint

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"main/pkg/types"

	"github.com/rs/zerolog"
)

const PaginationLimit = 1000

type RPC struct {
	URLs   []string
	Logger zerolog.Logger
}

func NewRPC(urls []string, logger zerolog.Logger) *RPC {
	return &RPC{
		URLs:   urls,
		Logger: logger.With().Str("component", "rpc").Logger(),
	}
}

func (rpc *RPC) GetAllProposals() ([]types.Proposal, error) {
	proposals := []types.Proposal{}
	offset := 0

	for {
		url := fmt.Sprintf(
			// 2 is for PROPOSAL_STATUS_VOTING_PERIOD
			"/cosmos/gov/v1beta1/proposals?pagination.limit=%d&pagination.offset=%d&proposal_status=2",
			PaginationLimit,
			offset,
		)

		var batchProposals types.ProposalsRPCResponse
		if err := rpc.Get(url, &batchProposals); err != nil {
			return nil, err
		}

		proposals = append(proposals, batchProposals.Proposals...)
		if len(batchProposals.Proposals) < PaginationLimit {
			break
		}

		offset += PaginationLimit
	}

	return proposals, nil
}

func (rpc *RPC) GetVote(proposal, voter string) (*types.VoteRPCResponse, error) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/votes/%s",
		proposal,
		voter,
	)

	var vote types.VoteRPCResponse
	if err := rpc.Get(url, &vote); err != nil {
		return nil, err
	}

	if vote.IsError() && !strings.Contains(vote.Message, "not found") {
		return nil, errors.New(vote.Message)
	}

	return &vote, nil
}

func (rpc *RPC) Get(url string, target interface{}) error {
	errors := make([]error, len(rpc.URLs))

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
		errors[index] = err
	}

	rpc.Logger.Warn().Str("url", url).Msg("All LCD requests failed")

	var sb strings.Builder

	sb.WriteString("All LCD requests failed:\n")
	for index, url := range rpc.URLs {
		sb.WriteString(fmt.Sprintf("#%d: %s -> %s\n", index+1, url, errors[index]))
	}

	return fmt.Errorf(sb.String())
}

func (rpc *RPC) GetFull(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * 1000000000}
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
