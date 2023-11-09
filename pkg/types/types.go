package types

import (
	"fmt"
	"main/pkg/utils"
	"time"

	"cosmossdk.io/math"
)

type Link struct {
	Name string
	Href string
}

func (l Link) Serialize() string {
	if l.Href == "" {
		return l.Name
	}

	return fmt.Sprintf("<a href='%s'>%s</a>", l.Href, l.Name)
}

type Proposal struct {
	ID          string
	Title       string
	Description string
	EndTime     time.Time
}

func (p Proposal) GetTimeLeft() string {
	return utils.FormatDuration(time.Until(p.EndTime).Round(time.Second))
}

func (p Proposal) GetProposalTime() string {
	return p.EndTime.Format(time.RFC1123)
}

type TallyInfo struct {
	Proposal Proposal
	Tally    Tally
	Pool     Pool
}

type ChainTallyInfos struct {
	Chain      *Chain
	TallyInfos []TallyInfo
}

func (t *TallyInfo) GetNotAbstained() math.LegacyDec {
	return t.Tally.Yes.Add(t.Tally.No).Add(t.Tally.NoWithVeto)
}

func (t *TallyInfo) GetTotalVoted() math.LegacyDec {
	return t.GetNotAbstained().Add(t.Tally.Abstain)
}

func (t *TallyInfo) GetQuorum() string {
	return fmt.Sprintf(
		"%.2f%%",
		t.GetTotalVoted().Quo(t.Pool.BondedTokens).Mul(math.LegacyNewDec(100)).MustFloat64(),
	)
}

func (t *TallyInfo) GetNotVoted() string {
	return fmt.Sprintf(
		"%.2f%%",
		math.LegacyNewDec(100).Sub(t.GetTotalVoted().Quo(t.Pool.BondedTokens).Mul(math.LegacyNewDec(100))).MustFloat64(),
	)
}

func (t *TallyInfo) GetAbstained() string {
	abstainedPercent := t.Tally.Abstain.Quo(t.GetTotalVoted()).Mul(math.LegacyNewDec(100)).MustFloat64()

	return fmt.Sprintf(
		"%.2f%%",
		abstainedPercent,
	)
}

func (t *TallyInfo) GetYesVotes() string {
	percent := t.Tally.Yes.Quo(t.GetTotalVoted()).Mul(math.LegacyNewDec(100)).MustFloat64()

	return fmt.Sprintf(
		"%.2f%%",
		percent,
	)
}

func (t *TallyInfo) GetNoVotes() string {
	percent := t.Tally.No.Quo(t.GetTotalVoted()).Mul(math.LegacyNewDec(100)).MustFloat64()

	return fmt.Sprintf(
		"%.2f%%",
		percent,
	)
}

func (t *TallyInfo) GetNoWithVetoVotes() string {
	percent := t.Tally.NoWithVeto.Quo(t.GetTotalVoted()).Mul(math.LegacyNewDec(100)).MustFloat64()

	return fmt.Sprintf(
		"%.2f%%",
		percent,
	)
}
