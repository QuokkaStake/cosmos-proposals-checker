package responses

import (
	"main/pkg/types"
	"time"
)

type ParamsResponse struct {
	Data struct {
		Threshold struct {
			ThresholdQuorum struct {
				Threshold struct {
					Percent float64 `json:"percent,string"`
				} `json:"threshold"`
				Quorum struct {
					Percent float64 `json:"percent,string"`
				} `json:"quorum"`
			} `json:"threshold_quorum"`
		} `json:"threshold"`
		MaxVotingPeriod struct {
			Time int `json:"time"`
		} `json:"max_voting_period"`
		AllowRevoting bool `json:"allow_revoting"`
	} `json:"data"`
}

func (params ParamsResponse) ToParams(chain *types.Chain) *types.ChainWithVotingParams {
	return &types.ChainWithVotingParams{
		Chain: chain,
		Params: []types.ChainParam{
			types.PercentParam{Description: "Threshold percent", Value: params.Data.Threshold.ThresholdQuorum.Threshold.Percent},
			types.PercentParam{Description: "Quorum percent", Value: params.Data.Threshold.ThresholdQuorum.Quorum.Percent},
			types.DurationParam{
				Description: "Max voting period",
				Value:       time.Duration(params.Data.MaxVotingPeriod.Time * 1e9),
			},
			types.BoolParam{Description: "Allow revoting", Value: params.Data.AllowRevoting},
		},
	}
}
