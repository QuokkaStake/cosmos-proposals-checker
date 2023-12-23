package responses

import (
	"main/pkg/types"
	"strconv"
	"time"
)

type ParamsResponse struct {
	Data struct {
		Threshold struct {
			ThresholdQuorum struct {
				Threshold struct {
					Percent string `json:"percent"`
				} `json:"threshold"`
				Quorum struct {
					Percent string `json:"percent"`
				} `json:"quorum"`
			} `json:"threshold_quorum"`
		} `json:"threshold"`
		MaxVotingPeriod struct {
			Time int `json:"time"`
		} `json:"max_voting_period"`
		AllowRevoting bool `json:"allow_revoting"`
	} `json:"data"`
}

func (params ParamsResponse) ToParams(chain *types.Chain) (*types.ChainWithVotingParams, []error) {
	thresholdPercent, err := strconv.ParseFloat(params.Data.Threshold.ThresholdQuorum.Threshold.Percent, 64)
	if err != nil {
		return nil, []error{err}
	}

	quorumPercent, err := strconv.ParseFloat(params.Data.Threshold.ThresholdQuorum.Quorum.Percent, 64)
	if err != nil {
		return nil, []error{err}
	}

	return &types.ChainWithVotingParams{
		Chain: chain,
		Params: []types.ChainParam{
			types.PercentParam{Description: "Threshold percent", Value: thresholdPercent},
			types.PercentParam{Description: "Quorum percent", Value: quorumPercent},
			types.DurationParam{
				Description: "Max voting period",
				Value:       time.Duration(params.Data.MaxVotingPeriod.Time * 1e9),
			},
			types.BoolParam{Description: "Allow revoting", Value: params.Data.AllowRevoting},
		},
	}, []error{}
}
