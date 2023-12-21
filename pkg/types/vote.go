package types

import (
	"main/pkg/utils"
	"strings"
)

// cosmos/gov/v1beta1/proposals/:id/votes/:wallet

type VoteOption struct {
	Option string
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
		return v.Option
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
