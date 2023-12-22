package responses

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"
	"time"
)

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

func (params ParamsResponse) ToParams(chain *types.Chain) (*types.ChainWithVotingParams, []error) {
	quorum, err := strconv.ParseFloat(params.TallyParams.Quorum, 64)
	if err != nil {
		return nil, []error{err}
	}

	threshold, err := strconv.ParseFloat(params.TallyParams.Threshold, 64)
	if err != nil {
		return nil, []error{err}
	}

	vetoThreshold, err := strconv.ParseFloat(params.TallyParams.VetoThreshold, 64)
	if err != nil {
		return nil, []error{err}
	}

	votingPeriod, err := time.ParseDuration(params.VotingParams.VotingPeriod)
	if err != nil {
		return nil, []error{err}
	}

	maxDepositPeriod, err := time.ParseDuration(params.DepositParams.MaxDepositPeriod)
	if err != nil {
		return nil, []error{err}
	}

	return &types.ChainWithVotingParams{
		Chain: chain,
		Params: []types.ChainParam{
			types.DurationParam{Description: "Voting period", Value: votingPeriod},
			types.DurationParam{Description: "Max deposit period", Value: maxDepositPeriod},
			types.AmountsParam{
				Description: "Min deposit amount",
				Value: utils.Map(params.DepositParams.MinDepositAmount, func(amount Amount) types.Amount {
					return types.Amount{
						Denom:  amount.Denom,
						Amount: amount.Amount,
					}
				}),
			},
			types.PercentParam{Description: "Quorum", Value: quorum},
			types.PercentParam{Description: "Threshold", Value: threshold},
			types.PercentParam{Description: "Veto threshold", Value: vetoThreshold},
		},
	}, nil
}
