package report

import (
	"main/pkg/events"
	"main/pkg/fs"
	"main/pkg/logger"
	mutes "main/pkg/mutes"
	"main/pkg/report/entry"
	reportersPkg "main/pkg/reporters"
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
	})

	err := dispatcher.Init()
	require.Error(t, err)
}

func TestReportDispatcherInitOk(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	})

	err := dispatcher.Init()
	require.NoError(t, err)
}

func TestReportDispatcherSendEmptyReport(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	})

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: make([]entry.ReportEntry, 0)})
}

func TestReportDispatcherSendReportDisabledReporter(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{WithDisabled: true},
	})

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}})
}

func TestReportDispatcherSendReportMuted(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	})

	err := dispatcher.Init()
	require.NoError(t, err)

	mutesManager.AddMute(&mutes.Mute{Expires: time.Now().Add(time.Minute)})

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.NotVotedEvent{
			Chain:    &types.Chain{Name: "chain"},
			Proposal: types.Proposal{ID: "proposal"},
		},
	}})
}

func TestReportDispatcherSendReportErrorSending(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{WithErrorSending: true},
	})

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}})
}

func TestReportDispatcherSendReportOk(t *testing.T) {
	t.Parallel()

	mutesManager := mutes.NewMutesManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())
	dispatcher := NewDispatcher(logger.GetNopLogger(), mutesManager, []reportersPkg.Reporter{
		&reportersPkg.TestReporter{},
	})

	err := dispatcher.Init()
	require.NoError(t, err)

	dispatcher.SendReport(reportersPkg.Report{Entries: []entry.ReportEntry{
		events.ProposalsQueryErrorEvent{},
	}})
}
