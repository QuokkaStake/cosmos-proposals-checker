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
func TestTelegramReporterListProposalsOk(t *testing.T) {
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
			Text:   "/proposals",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = database.UpsertMute(&types.Mute{})
	require.NoError(t, err)

	err = reporter.HandleProposals(ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterListProposalsRenderOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/list-proposals.html")),
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
			Text:   "/proposals",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	timeParsed, err := time.Parse(time.RFC3339, "2024-12-03T10:13:01Z")
	require.NoError(t, err)

	renderTime, err := time.Parse(time.RFC3339, "2024-12-01T16:56:01Z")
	require.NoError(t, err)

	err = reporter.ReplyRender(ctx, "proposals", state.RenderedState{
		RenderTime: renderTime,
		ChainInfos: []state.RenderedChainInfo{
			{
				Chain:          &types.Chain{Name: "chain1"},
				ProposalsError: &types.QueryError{QueryError: errors.New("proposals fetch error")},
			},
			{
				Chain: &types.Chain{
					Name:       "chain2",
					PrettyName: "FancyChainName",
				},
				ProposalVotes: []state.RenderedProposalVotes{
					{
						Proposal: types.Proposal{
							ID:          "proposal1",
							Title:       "proposal1title",
							Description: "proposal1description",
							EndTime:     timeParsed,
							Status:      "PROPOSAL_STATUS_VOTING_PERIOD",
						},
						Votes: []state.RenderedWalletVote{
							{
								Wallet: &types.Wallet{Address: "wallet1"},
								Error:  &types.QueryError{QueryError: errors.New("vote fetch error")},
							},
							{
								Vote:   nil,
								Wallet: &types.Wallet{Address: "wallet2"},
							},
							{
								Wallet: &types.Wallet{Address: "wallet3", Alias: "FancyWalletAlias"},
								Vote: &types.Vote{
									ProposalID: "proposal1",
									Voter:      "wallet3",
									Options: types.VoteOptions{
										{Option: "Yes", Weight: 1},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
}
