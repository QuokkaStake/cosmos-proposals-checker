package responses

import (
	"main/pkg/types"
)

type Vote struct {
	Voter string `json:"voter"`
	Vote  string `json:"vote"`
}

func (v Vote) GetOption() string {
	options := map[string]string{
		"yes":     "👌Yes",
		"no":      "🚫No",
		"abstain": "🤷Abstain",
	}

	if option, ok := options[v.Vote]; ok {
		return option
	}

	return v.Vote
}

type VoteResponse struct {
	Data struct {
		Vote *Vote `json:"vote"`
	} `json:"data"`
}

func (v VoteResponse) ToVote(proposalID string) *types.Vote {
	if v.Data.Vote == nil {
		return nil
	}

	return &types.Vote{
		ProposalID: proposalID,
		Voter:      v.Data.Vote.Voter,
		Options: []types.VoteOption{
			{Option: v.Data.Vote.GetOption(), Weight: 1},
		},
	}
}
