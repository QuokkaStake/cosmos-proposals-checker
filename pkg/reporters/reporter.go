package reporters

import (
	"main/pkg/report/entry"
)

type Reporter interface {
	Init() error
	Enabled() bool
	SendReport(report Report) error
	Name() string
}

type Report struct {
	Entries []entry.ReportEntry
}

func (r *Report) Empty() bool {
	return len(r.Entries) == 0
}
