package types

import (
	"fmt"
	"main/pkg/utils"
	"time"

	"cosmossdk.io/math"
)

type TallyOption struct {
	Option string
	Voted  math.LegacyDec
}

type Tally []TallyOption

func (t Tally) GetTotalVoted() math.LegacyDec {
	sum := math.LegacyNewDec(0)

	for _, option := range t {
		sum = sum.Add(option.Voted)
	}

	return sum
}

func (t Tally) GetVoted(option TallyOption) string {
	votedPercent := option.Voted.
		Quo(t.GetTotalVoted()).
		Mul(math.LegacyNewDec(100)).
		MustFloat64()

	return fmt.Sprintf(
		"%.2f%%",
		votedPercent,
	)
}

type TallyInfo struct {
	Proposal         Proposal
	Tally            Tally
	TotalVotingPower math.LegacyDec
}

func (t TallyInfo) GetQuorum() string {
	return fmt.Sprintf(
		"%.2f%%",
		t.Tally.GetTotalVoted().
			Quo(t.TotalVotingPower).
			Mul(math.LegacyNewDec(100)).
			MustFloat64(),
	)
}

func (t TallyInfo) GetNotVoted() string {
	return fmt.Sprintf(
		"%.2f%%",
		math.LegacyNewDec(100).
			Sub(t.Tally.GetTotalVoted().
				Quo(t.TotalVotingPower).
				Mul(math.LegacyNewDec(100)),
			).MustFloat64(),
	)
}

type ChainsTallyInfos struct {
	RenderTime       time.Time
	ChainsTallyInfos map[string]ChainTallyInfos
}

type ChainTallyInfos struct {
	Chain      *Chain
	TallyInfos []TallyInfo
}

func (i ChainsTallyInfos) GetProposalTimeLeft(p Proposal) string {
	return utils.FormatDuration(p.EndTime.Sub(i.RenderTime).Round(time.Second))
}
