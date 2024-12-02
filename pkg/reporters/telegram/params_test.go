package telegram

import (
	"main/assets"
	"main/pkg/data"
	databasePkg "main/pkg/database"
	"main/pkg/fetchers"
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
func TestTelegramReporterGetParamsError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error getting chain params: got 3 errors when fetching chain params"),
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
			Text:   "/params",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = database.UpsertMute(&types.Mute{})
	require.NoError(t, err)

	err = reporter.HandleParams(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterGetParamsOk(t *testing.T) {
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
	logger := loggerPkg.GetDefaultLogger()
	tracer := tracing.InitNoopTracer()
	database := &databasePkg.StubDatabase{}
	mutesManager := mutesmanager.NewMutesManager(logger, database)
	stateGenerator := state.NewStateGenerator(logger, tracer, chains)
	dataManager := data.NewManager(logger, chains, tracer)
	dataManager.Fetchers = []fetchers.Fetcher{
		&fetchers.TestFetcher{},
	}

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
			Text:   "/params",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = database.UpsertMute(&types.Mute{})
	require.NoError(t, err)

	err = reporter.HandleParams(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterGetParamsRenderOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/params.html")),
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
			Text:   "/params",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.ReplyRender(ctx, "params", map[string]types.ChainWithVotingParams{
		"chain": {
			Chain: &types.Chain{Name: "chain", PrettyName: "FancyChainName"},
			Params: []types.ChainParam{
				types.PercentParam{Description: "percentParam", Value: 0.7},
				types.DurationParam{Description: "durationParam", Value: time.Hour},
				types.BoolParam{Description: "boolParam", Value: false},
				types.BoolParam{Description: "boolParamSecond", Value: true},
				types.AmountsParam{
					Description: "amountsParam",
					Value:       []types.Amount{{Amount: "100", Denom: "denom"}},
				},
			},
		},
	})
	require.NoError(t, err)
}
