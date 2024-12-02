package telegram

import (
	"context"
	"errors"
	"main/assets"
	"main/pkg/data"
	databasePkg "main/pkg/database"
	"main/pkg/events"
	loggerPkg "main/pkg/logger"
	mutesmanager "main/pkg/mutes"
	"main/pkg/report/entry"
	"main/pkg/state"
	"main/pkg/tracing"
	"main/pkg/types"
	"strings"
	"sync"
	"testing"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled
func TestReporterInitNotEnabled(t *testing.T) {
	config := types.TelegramConfig{}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)
	err = reporter.Init()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestReporterInitCannotFetchBot(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)
	err = reporter.Init()
	require.Error(t, err)
}

//nolint:paralleltest // disabled
func TestReporterInitOkay(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_ = reporter.Init()

		ctx := reporter.TelegramBot.NewContext(tele.Update{
			ID: 1,
			Message: &tele.Message{
				Sender: &tele.User{Username: "testuser"},
				Text:   "/help",
				Chat:   &tele.Chat{ID: 2},
			},
		})
		reporter.TelegramBot.OnError(errors.New("custom error"), ctx)
		reporter.Stop()
		wg.Done()
	}()

	wg.Wait()
}

//nolint:paralleltest // disabled
func TestAppBotSendMultilineFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewErrorResponder(errors.New("custom error")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	err = reporter.InitBot()
	require.NoError(t, err)

	ctx := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.BotReply(ctx, strings.Repeat("a", 5000))
	require.Error(t, err)

	err = reporter.BotReply(ctx, strings.Repeat("a", 10))
	require.Error(t, err)
}

//nolint:paralleltest // disabled
func TestAppBotSendMultilineOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	err = reporter.InitBot()
	require.NoError(t, err)

	ctx := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.BotReply(ctx, strings.Repeat("a", 5000))
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestAppBotSendReportEntryFailedToSerialize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	err = reporter.InitBot()
	require.NoError(t, err)

	err = reporter.SendReportEntry(&events.NotExistingEvent{}, context.Background())
	require.Error(t, err)
}

//nolint:paralleltest // disabled
func TestAppBotSendReportEntryFailedToSend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewErrorResponder(errors.New("custom error")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	err = reporter.InitBot()
	require.NoError(t, err)

	err = reporter.SendReportEntry(&events.ProposalsQueryErrorEvent{
		Chain: &types.Chain{Name: "chain"},
		Error: &types.QueryError{QueryError: errors.New("query error")},
	}, context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}

//nolint:paralleltest // disabled
func TestAppBotSendReportEntryOk(t *testing.T) {
	httpmock.Activate()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewTelegramReporter(
		config,
		mutesManager,
		stateGenerator,
		dataManager,
		logger,
		"1.2.3",
		timezone,
		tracer,
	)

	err = reporter.InitBot()
	require.NoError(t, err)

	httpmock.Reset()

	proposalEndTime, err := time.Parse(time.RFC3339, "2024-12-03T10:13:01Z")
	require.NoError(t, err)

	renderTime, err := time.Parse(time.RFC3339, "2024-12-01T16:56:01Z")
	require.NoError(t, err)

	inputs := []struct {
		event      entry.ReportEntry
		resultFile string
	}{
		{
			event: events.ProposalsQueryErrorEvent{
				Chain: &types.Chain{Name: "chain"},
				Error: &types.QueryError{QueryError: errors.New("query error")},
			},
			resultFile: "responses/telegram-proposals-query-error.html",
		},
		{
			event: events.GenericErrorEvent{
				Chain: &types.Chain{Name: "chain"},
				Error: &types.QueryError{QueryError: errors.New("query error")},
			},
			resultFile: "responses/telegram-generic-error.html",
		},
		{
			event: events.FinishedVotingEvent{
				Chain: &types.Chain{Name: "chain"},
				Proposal: types.Proposal{
					ID:     "proposal",
					Title:  "proposal title",
					Status: types.ProposalStatusPassed,
				},
			},
			resultFile: "responses/telegram-voting-finished.html",
		},
		{
			event: events.VoteQueryError{
				Chain: &types.Chain{Name: "chain"},
				Proposal: types.Proposal{
					ID: "proposal",
				},
				Error: &types.QueryError{QueryError: errors.New("query error")},
			},
			resultFile: "responses/telegram-vote-query-error.html",
		},
		{
			event: events.NotVotedEvent{
				RenderTime: renderTime,
				Chain:      &types.Chain{Name: "chain"},
				Wallet:     &types.Wallet{Address: "address"},
				Proposal: types.Proposal{
					ID:      "proposal",
					Title:   "proposal title",
					EndTime: proposalEndTime,
				},
			},
			resultFile: "responses/telegram-not-voted.html",
		},
		{
			event: events.VotedEvent{
				RenderTime: renderTime,
				Chain:      &types.Chain{Name: "chain"},
				Wallet:     &types.Wallet{Address: "address"},
				Vote: &types.Vote{
					Options: types.VoteOptions{
						{Option: "Yes", Weight: 1},
					},
				},
				Proposal: types.Proposal{
					ID:      "proposal",
					Title:   "proposal title",
					EndTime: proposalEndTime,
				},
			},
			resultFile: "responses/telegram-voted.html",
		},
		{
			event: events.RevotedEvent{
				RenderTime: renderTime,
				Chain:      &types.Chain{Name: "chain"},
				Wallet:     &types.Wallet{Address: "address"},
				Vote: &types.Vote{
					Options: types.VoteOptions{
						{Option: "Yes", Weight: 1},
					},
				},
				OldVote: &types.Vote{
					Options: types.VoteOptions{
						{Option: "No", Weight: 1},
					},
				},
				Proposal: types.Proposal{
					ID:      "proposal",
					Title:   "proposal title",
					EndTime: proposalEndTime,
				},
			},
			resultFile: "responses/telegram-revoted.html",
		},
	}

	for _, input := range inputs {
		t.Run(input.event.Name(), func(t *testing.T) {
			httpmock.RegisterMatcherResponder(
				"POST",
				"https://api.telegram.org/botxxx:yyy/sendMessage",
				types.TelegramResponseHasBytes(assets.GetBytesOrPanic(input.resultFile)),
				httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

			defer httpmock.Reset()

			err = reporter.SendReportEntry(input.event, context.Background())
			require.NoError(t, err)
		})
	}
}
