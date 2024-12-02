package telegram

import (
	"errors"
	"main/assets"
	"main/pkg/data"
	databasePkg "main/pkg/database"
	loggerPkg "main/pkg/logger"
	mutesmanager "main/pkg/mutes"
	"main/pkg/state"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/guregu/null/v5"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteNotEnoughArgs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error muting notification: Usage: /proposals_mute <duration> [params]"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteInvalidDuration(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error muting notification: Invalid duration provided: invalid"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute invalid",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteInvalidArgument(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error muting notification: Invalid param at position 1: expected an expression like \"[chain=cosmos]\", but got invalid"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute 48h invalid",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error adding mute: custom error"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	config := types.TelegramConfig{TelegramToken: "xxx:yyy", TelegramChat: 123}
	chains := types.Chains{{Name: "chain", LCDEndpoints: []string{"https://example.com"}}}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{UpsertMuteError: errors.New("custom error")}
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
			Text:   "/proposals_mute 48h",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute 48h",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterAddMuteOkWithParams(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute 48h chain=chain proposal=proposal",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.HandleAddMute(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterRenderMuteNoParams(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/add-mutes-no-params.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute 48h",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	timeParsed, err := time.Parse(time.RFC3339, "2024-12-01T16:56:01Z")
	require.NoError(t, err)

	mute := types.Mute{
		Expires: timeParsed,
	}

	err = reporter.ReplyRender(ctx, "mute_added", mute)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReporterRenderMuteWithParams(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/add-mutes-params.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

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
			Text:   "/proposals_mute 48h",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	timeParsed, err := time.Parse(time.RFC3339, "2024-12-01T16:56:01Z")
	require.NoError(t, err)

	mute := types.Mute{
		Chain:      null.StringFrom("chain"),
		ProposalID: null.StringFrom("proposal"),
		Comment:    "comment",
		Expires:    timeParsed,
	}

	err = reporter.ReplyRender(ctx, "mute_added", mute)
	require.NoError(t, err)
}
