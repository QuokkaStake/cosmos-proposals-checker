package reporters

import (
	"errors"
	"main/pkg/report/entry"
)

type TestReporter struct {
	WithInitFail     bool
	WithDisabled     bool
	WithErrorSending bool
}

func (r *TestReporter) Init() error {
	if r.WithInitFail {
		return errors.New("fail")
	}

	return nil
}

func (r *TestReporter) Enabled() bool {
	return !r.WithDisabled
}

func (r *TestReporter) SendReportEntry(entry entry.ReportEntry) error {
	if r.WithErrorSending {
		return errors.New("fail")
	}

	return nil
}

func (r *TestReporter) Name() string {
	return "test-reporter"
}
