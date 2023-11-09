package types

import (
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

// cosmos/gov/v1beta1/proposals/:id/votes/:wallet

type Vote struct {
	ProposalID string       `json:"proposal_id"`
	Voter      string       `json:"voter"`
	Option     string       `json:"option"`
	Options    []VoteOption `json:"options"`
}

type VoteOption struct {
	Option string `json:"option"`
	Weight string `json:"weight"`
}

func (v Vote) ResolveVote() string {
	if len(v.Options) > 0 {
		optionsStrings := utils.Map(v.Options, func(v VoteOption) string {
			return utils.ResolveVote(v.Option)
		})

		return strings.Join(optionsStrings, ", ")
	}

	return utils.ResolveVote(v.Option)
}

type VoteRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Vote    *Vote  `json:"vote"`
}

func (v VoteRPCResponse) IsError() bool {
	return v.Code != 0
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
