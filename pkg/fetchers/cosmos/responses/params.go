package responses

import (
	"main/pkg/types"
	"main/pkg/utils"

	"cosmossdk.io/math"
)

type ParamsResponse struct {
	VotingParams  VotingParams  `json:"voting_params"`
	DepositParams DepositParams `json:"deposit_params"`
	TallyParams   TallyParams   `json:"tally_params"`
}

type VotingParams struct {
	VotingPeriod types.Duration `json:"voting_period"`
}

type DepositParams struct {
	MinDepositAmount []Amount       `json:"min_deposit"`
	MaxDepositPeriod types.Duration `json:"max_deposit_period"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type TallyParams struct {
	Quorum        math.LegacyDec `json:"quorum"`
	Threshold     math.LegacyDec `json:"threshold"`
	VetoThreshold math.LegacyDec `json:"veto_threshold"`
}

func (params ParamsResponse) ToParams(chain *types.Chain) (*types.ChainWithVotingParams, []error) {
	return &types.ChainWithVotingParams{
		Chain: chain,
		Params: []types.ChainParam{
			types.DurationParam{Description: "Voting period", Value: params.VotingParams.VotingPeriod.Duration},
			types.DurationParam{Description: "Max deposit period", Value: params.DepositParams.MaxDepositPeriod.Duration},
			types.AmountsParam{
				Description: "Min deposit amount",
				Value: utils.Map(params.DepositParams.MinDepositAmount, func(amount Amount) types.Amount {
					return types.Amount{
						Denom:  amount.Denom,
						Amount: amount.Amount,
					}
				}),
			},
			types.PercentParam{Description: "Quorum", Value: params.TallyParams.Quorum.MustFloat64()},
			types.PercentParam{Description: "Threshold", Value: params.TallyParams.Threshold.MustFloat64()},
			types.PercentParam{Description: "Veto threshold", Value: params.TallyParams.VetoThreshold.MustFloat64()},
		},
	}, nil
}
