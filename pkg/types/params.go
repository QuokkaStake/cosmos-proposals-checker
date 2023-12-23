package types

import (
	"fmt"
	"strings"
	"time"

	"main/pkg/utils"
)

type ChainWithVotingParams struct {
	Chain  *Chain
	Params []ChainParam
}

type ChainParam interface {
	GetDescription() string
	Serialize() string
}

/* ------------------------------------------------ */

type PercentParam struct {
	Description string
	Value       float64
}

func (p PercentParam) GetDescription() string {
	return p.Description
}

func (p PercentParam) Serialize() string {
	return fmt.Sprintf("%.2f%%", p.Value*100)
}

/* ------------------------------------------------ */

type DurationParam struct {
	Description string
	Value       time.Duration
}

func (p DurationParam) GetDescription() string {
	return p.Description
}

func (p DurationParam) Serialize() string {
	return utils.FormatDuration(p.Value)
}

/* ------------------------------------------------ */

type BoolParam struct {
	Description string
	Value       bool
}

func (p BoolParam) GetDescription() string {
	return p.Description
}

func (p BoolParam) Serialize() string {
	if p.Value {
		return "Yes"
	}

	return "No"
}

/* ------------------------------------------------ */

type Amount struct {
	Denom  string
	Amount string
}

type AmountsParam struct {
	Description string
	Value       []Amount
}

func (p AmountsParam) GetDescription() string {
	return p.Description
}

func (p AmountsParam) Serialize() string {
	amountsAsStrings := utils.Map(p.Value, func(a Amount) string {
		return a.Amount + " " + a.Denom
	})
	return strings.Join(amountsAsStrings, ",")
}
