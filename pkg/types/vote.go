package types

import (
	"main/pkg/utils"
	"strings"
)

// cosmos/gov/v1beta1/proposals/:id/votes/:wallet

type VoteType string

func (v VoteType) Resolve() string {
	votes := map[string]string{
		"VOTE_OPTION_YES":          "Yes",
		"VOTE_OPTION_ABSTAIN":      "Abstain",
		"VOTE_OPTION_NO":           "No",
		"VOTE_OPTION_NO_WITH_VETO": "No with veto",
	}

	if vote, ok := votes[string(v)]; ok && v != "" {
		return vote
	}

	return string(v)
}

type VoteOption struct {
	Option VoteType
	Weight float64
}
type VoteOptions []VoteOption

type Vote struct {
	ProposalID string
	Voter      string
	Options    VoteOptions
}

func (v Vote) ResolveVote() string {
	optionsStrings := utils.Map(v.Options, func(v VoteOption) string {
		return v.Option.Resolve()
	})

	return strings.Join(optionsStrings, ", ")
}

func (v Vote) VotesEquals(other *Vote) bool {
	if len(v.Options) != len(other.Options) {
		return false
	}

	for index, option := range v.Options {
		if option.Weight != other.Options[index].Weight || option.Option != other.Options[index].Option {
			return false
		}
	}

	return true
}
