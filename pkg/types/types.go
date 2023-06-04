package types

import (
	"fmt"
	"main/pkg/utils"
	"time"
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
