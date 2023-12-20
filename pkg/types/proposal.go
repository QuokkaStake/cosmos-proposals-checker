package types

import (
	"main/pkg/utils"
	"time"
)

type Proposal struct {
	ID          string
	Title       string
	Description string
	EndTime     time.Time
}

func (p Proposal) GetTimeLeft() string {
	return utils.FormatDuration(time.Until(p.EndTime).Round(time.Second))
}
