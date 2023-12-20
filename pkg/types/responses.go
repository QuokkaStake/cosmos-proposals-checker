package types

import (
	"fmt"
	"strings"
	"time"

	"main/pkg/utils"

	"cosmossdk.io/math"
)

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
