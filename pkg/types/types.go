package types

import (
	"fmt"
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

type ReportEntryType string

const (
	NotVoted           ReportEntryType = "not_voted"
	Voted              ReportEntryType = "voted"
	Revoted            ReportEntryType = "revoted"
	ProposalQueryError ReportEntryType = "proposal_query_error"
	VoteQueryError     ReportEntryType = "vote_query_error"
)
