package reporters

import (
	"main/pkg/report/entry"
)

type Reporter interface {
	Init() error
	Enabled() bool
	SendReportEntry(entry entry.ReportEntry) error
	Name() string
}

type Report struct {
	Entries []entry.ReportEntry
}

func (r *Report) Empty() bool {
	return len(r.Entries) == 0
}
