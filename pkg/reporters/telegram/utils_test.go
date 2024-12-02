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

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestTelegramReporterReplyRenderFailedToRender(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error rendering template: template: pattern matches no files: `telegram/not_found.html`"), //nolint:dupword
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
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.ReplyRender(ctx, "not_found", nil)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterReplyRenderFailedToSend(t *testing.T) {
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

	err = reporter.ReplyRender(ctx, "help", "1.2.3")
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}

//nolint:paralleltest // disabled
func TestTelegramReporterEditRenderFailedToRender(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error rendering template: template: pattern matches no files: `telegram/not_found.html`"), //nolint:dupword
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
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.EditRender(ctx, ctx.Message(), "not_found", nil)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterEditRenderFailedToSend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/editMessageText",
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

	err = reporter.EditRender(ctx, ctx.Message(), "help", "1.2.3")
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}
