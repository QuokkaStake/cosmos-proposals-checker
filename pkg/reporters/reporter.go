package reporters

import (
	"time"

	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/utils"
)

type Reporter interface {
	Init()
	Enabled() bool
	SendReport(report Report) error
	Name() string
}

type Report struct {
	Entries []ReportEntry
}

func (r *Report) Empty() bool {
	return len(r.Entries) == 0
}

type ReportEntry struct {
	Chain                  *configTypes.Chain
	Wallet                 *configTypes.Wallet
	ProposalID             string
	ProposalTitle          string
	ProposalDescription    string
	ProposalVoteEndingTime time.Time
	Type                   types.ReportEntryType
	Value                  string
	OldValue               string
}

func (e ReportEntry) HasVoted() bool {
	return e.Value != ""
}

func (e ReportEntry) HasRevoted() bool {
	return e.Value != "" && e.OldValue != ""
}

func (e ReportEntry) IsVoteOrNotVoted() bool {
	return e.Type == types.NotVoted || e.Type == types.Voted
}

func (e ReportEntry) GetProposalTime() string {
	return e.ProposalVoteEndingTime.Format(time.RFC1123)
}

func (e ReportEntry) GetProposalTimeLeft() string {
	return utils.FormatDuration(time.Until(e.ProposalVoteEndingTime).Round(time.Second))
}

func (e ReportEntry) GetVote() string {
	return utils.ResolveVote(e.Value)
}

func (e ReportEntry) GetOldVote() string {
	return utils.ResolveVote(e.OldValue)
}
