package report

import (
	"context"
	"main/pkg/events"
	"main/pkg/fs"
	"main/pkg/logger"
	mutes "main/pkg/mutes"
	"main/pkg/report/entry"
	reportersPkg "main/pkg/reporters"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReportDispatcherInitFail(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{WithInitFail: true},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.Error(t, err)
}

func TestReportDispatcherInitOk(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)
}

func TestReportDispatcherSendEmptyReport(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: make([]entry.ReportEntry, 0)}, context.Background())
}

func TestReportDispatcherSendReportDisabledReporter(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{WithDisabled: true},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}}, context.Background())
}

func TestReportDispatcherSendReportMuted(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)

	mutesManager.AddMute(&mutes.Mute{Expires: time.Now().Add(time.Minute)})

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.NotVotedEvent{
			Chain:    &types.Chain{Name: "chain"},
			Proposal: types.Proposal{ID: "proposal"},
		},
	}}, context.Background())
}

func TestReportDispatcherSendReportErrorSending(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{WithErrorSending: true},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}}, context.Background())
}

func TestReportDispatcherSendReportOk(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	}, tracing.InitNoopTracer())

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}}, context.Background())
}
