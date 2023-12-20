package responses

import (
	"main/pkg/types"
	"strconv"
)

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

type VoteRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Vote    *Vote  `json:"vote"`
}

func (v VoteRPCResponse) IsError() bool {
	return v.Code != 0
}

func (v VoteRPCResponse) ToVote() (*types.Vote, error) {
	var options []types.VoteOption

	if len(v.Vote.Options) > 0 {
		options = make([]types.VoteOption, len(v.Vote.Options))

		for index, option := range v.Vote.Options {
			weight, err := strconv.ParseFloat(option.Weight, 64)
			if err != nil {
				return nil, err
			}

			options[index] = types.VoteOption{
				Option: types.VoteType(option.Option),
				Weight: weight,
			}
		}
	} else {
		options = make([]types.VoteOption, 1)
		options[0] = types.VoteOption{
			Option: types.VoteType(v.Vote.Option),
			Weight: 1,
		}
	}

	return &types.Vote{
		ProposalID: v.Vote.ProposalID,
		Voter:      v.Vote.Voter,
		Options:    options,
	}, nil
}
