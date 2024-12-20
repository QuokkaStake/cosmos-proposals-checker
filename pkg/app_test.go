package pkg

import (
	databasePkg "main/pkg/database"
	"main/pkg/fs"
	reportersPkg "main/pkg/reporters"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/require"
)

func TestAppFailToLoadConfig(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	filesystem := &fs.TestFS{}

	NewApp("notexisting.toml", filesystem, "1.2.3")
}

func TestAppFailInvalidConfig(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	filesystem := &fs.TestFS{}

	NewApp("config-invalid.toml", filesystem, "1.2.3")
}

func TestAppCreateConfigWithWarnings(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	app := NewApp("config-with-warnings.toml", filesystem, "1.2.3")
	require.NotNil(t, app)
}

func TestAppCreateConfigValid(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	app := NewApp("config-valid.toml", filesystem, "1.2.3")
	require.NotNil(t, app)
}

func TestAppStartReporterFailedToInit(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	app := NewApp("config-valid.toml", filesystem, "1.2.3")
	app.ReportDispatcher.Reporters = []reportersPkg.Reporter{&reportersPkg.TestReporter{WithInitFail: true}}

	app.Start()
	require.NotNil(t, app)
}

func TestAppStartInvalidCronPattern(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	app := NewApp("config-valid.toml", filesystem, "1.2.3")
	app.Config.Interval = "invalid"
	app.Database = &databasePkg.StubDatabase{}
	app.ReportGenerator.Database = &databasePkg.StubDatabase{}

	app.Start()
	require.NotNil(t, app)
}

func TestAppStartOk(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}
	app := NewApp("config-valid.toml", filesystem, "1.2.3")
	app.ReportDispatcher.Reporters = []reportersPkg.Reporter{&reportersPkg.TestReporter{}}
	app.Database = &databasePkg.StubDatabase{}
	app.ReportGenerator.Database = &databasePkg.StubDatabase{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		app.Start()
		wg.Done()
	}()

	app.Stop()
	wg.Wait()
}

//nolint:paralleltest // disabled due to httpmock usage
func TestAppReport(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	filesystem := &fs.TestFS{}
	app := NewApp("config-valid.toml", filesystem, "1.2.3")
	app.Database = &databasePkg.StubDatabase{}
	app.ReportGenerator.Database = &databasePkg.StubDatabase{}
	app.Report()
}
