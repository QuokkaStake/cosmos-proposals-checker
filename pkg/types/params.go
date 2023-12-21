package types

import (
	"fmt"
	"strings"
	"time"

	"main/pkg/utils"
)

type Amount struct {
	Denom  string
	Amount string
}

type ChainWithVotingParams struct {
	Chain            *Chain
	VotingPeriod     time.Duration
	MinDepositAmount []Amount
	MaxDepositPeriod time.Duration
	Quorum           float64
	Threshold        float64
	VetoThreshold    float64
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
