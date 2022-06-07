package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

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

func (rpc *RPC) GetAllProposals() ([]Proposal, error) {
	proposals := []Proposal{}
	offset := 0

	for {
		url := fmt.Sprintf(
			// 2 is for PROPOSAL_STATUS_VOTING_PERIOD
			"/cosmos/gov/v1beta1/proposals?pagination.limit=%d&pagination.offset=%d&proposal_status=2",
			PaginationLimit,
			offset,
		)

		var batchProposals ProposalsRPCResponse
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

func (rpc *RPC) GetVote(proposal, voter string) (*VoteRPCResponse, error) {
	url := fmt.Sprintf(
		"/cosmos/gov/v1beta1/proposals/%s/votes/%s",
		proposal,
		voter,
	)

	var vote VoteRPCResponse
	if err := rpc.Get(url, &vote); err != nil {
		return nil, err
	}

	return &vote, nil
}

func (rpc *RPC) Get(url string, target interface{}) error {
	for _, lcd := range rpc.URLs {
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
	}

	rpc.Logger.Warn().Str("url", url).Msg("All LCD requests failed")
	return fmt.Errorf("all LCD requests failed")
}

func (rpc *RPC) GetFull(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * 1000000000}
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

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
