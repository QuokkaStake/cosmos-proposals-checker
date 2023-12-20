package types

import (
	"fmt"
	"strings"
	"time"

	"main/pkg/utils"

	"cosmossdk.io/math"
)

// cosmos/gov/v1beta1/proposals?pagination.limit=1000&pagination.offset=0

type V1beta1Proposal struct {
	ProposalID    string           `json:"proposal_id"`
	Status        string           `json:"status"`
	Content       *ProposalContent `json:"content"`
	VotingEndTime time.Time        `json:"voting_end_time"`
}

func (p V1beta1Proposal) ToProposal() Proposal {
	return Proposal{
		ID:          p.ProposalID,
		Title:       p.Content.Title,
		Description: p.Content.Description,
		EndTime:     p.VotingEndTime,
	}
}

type ProposalContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type V1Beta1ProposalsRPCResponse struct {
	Code      int64             `json:"code"`
	Message   string            `json:"message"`
	Proposals []V1beta1Proposal `json:"proposals"`
}

// cosmos/gov/v1beta1/proposals?pagination.limit=1000&pagination.offset=0

type V1ProposalMessage struct {
	Content ProposalContent `json:"content"`
}

type V1Proposal struct {
	ProposalID    string              `json:"id"`
	Status        string              `json:"status"`
	VotingEndTime time.Time           `json:"voting_end_time"`
	Messages      []V1ProposalMessage `json:"messages"`

	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func (p V1Proposal) ToProposal() Proposal {
	// Some chains (namely, Quicksilver) do not have title and description fields,
	// instead they have content.title and content.description per each message.
	// Others (namely, Kujira) have title and summary text.
	// This should work for all of them.
	title := p.Title
	if title == "" {
		titles := utils.Map(p.Messages, func(m V1ProposalMessage) string {
			return m.Content.Title
		})

		title = strings.Join(titles, ", ")
	}

	description := p.Summary
	if description == "" {
		descriptions := utils.Map(p.Messages, func(m V1ProposalMessage) string {
			return m.Content.Description
		})

		description = strings.Join(descriptions, ", ")
	}

	return Proposal{
		ID:          p.ProposalID,
		Title:       title,
		Description: description,
		EndTime:     p.VotingEndTime,
	}
}

type V1ProposalsRPCResponse struct {
	Code      int64        `json:"code"`
	Message   string       `json:"message"`
	Proposals []V1Proposal `json:"proposals"`
}

type TallyRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Tally   *Tally `json:"tally"`
}

type Tally struct {
	Yes        math.LegacyDec `json:"yes"`
	No         math.LegacyDec `json:"no"`
	NoWithVeto math.LegacyDec `json:"no_with_veto"`
	Abstain    math.LegacyDec `json:"abstain"`
}

type PoolRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Pool    *Pool  `json:"pool"`
}

type Pool struct {
	BondedTokens math.LegacyDec `json:"bonded_tokens"`
}

type ParamsResponse struct {
	VotingParams  VotingParams  `json:"voting_params"`
	DepositParams DepositParams `json:"deposit_params"`
	TallyParams   TallyParams   `json:"tally_params"`
}

type VotingParams struct {
	VotingPeriod string `json:"voting_period"`
}

type DepositParams struct {
	MinDepositAmount []Amount `json:"min_deposit"`
	MaxDepositPeriod string   `json:"max_deposit_period"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type TallyParams struct {
	Quorum        string `json:"quorum"`
	Threshold     string `json:"threshold"`
	VetoThreshold string `json:"veto_threshold"`
}

type ChainWithVotingParams struct {
	Chain            *Chain
	VotingPeriod     time.Duration `json:"voting_period"`
	MinDepositAmount []Amount      `json:"amount"`
	MaxDepositPeriod time.Duration `json:"max_deposit_period"`
	Quorum           float64       `json:"quorum"`
	Threshold        float64       `json:"threshold"`
	VetoThreshold    float64       `json:"veto_threshold"`
}

func (c ChainWithVotingParams) FormatQuorum() string {
	return fmt.Sprintf("%.2f%%", c.Quorum*100)
}

func (c ChainWithVotingParams) FormatThreshold() string {
	return fmt.Sprintf("%.2f%%", c.Threshold*100)
}

func (c ChainWithVotingParams) FormatVetoThreshold() string {
	return fmt.Sprintf("%.2f%%", c.VetoThreshold*100)
}

func (c ChainWithVotingParams) FormatMinDepositAmount() string {
	amountsAsStrings := utils.Map(c.MinDepositAmount, func(a Amount) string {
		return a.Amount + " " + a.Denom
	})
	return strings.Join(amountsAsStrings, ",")
}
