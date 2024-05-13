package pkg

import (
	"main/pkg/data"
	"main/pkg/fs"
	"main/pkg/logger"
	mutes "main/pkg/mutes"
	"main/pkg/report"
	reportersPkg "main/pkg/reporters"
	"main/pkg/reporters/pagerduty"
	"main/pkg/reporters/telegram"
	"main/pkg/state"
	"main/pkg/types"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type App struct {
	Logger           *zerolog.Logger
	Config           *types.Config
	StateManager     *state.Manager
	ReportGenerator  *report.Generator
	StateGenerator   *state.Generator
	ReportDispatcher *report.Dispatcher
}

func NewApp(configPath string, filesystem fs.FS, version string) *App {
	config, err := GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	if warnings := config.DisplayWarnings(); len(warnings) > 0 {
		config.LogWarnings(logger.GetDefaultLogger(), warnings)
	} else {
		logger.GetDefaultLogger().Info().Msg("Provided config is valid.")
	}

	log := logger.GetLogger(config.LogConfig)

	stateManager := state.NewStateManager(config.StatePath, filesystem, log)
	mutesManager := mutes.NewMutesManager(config.MutesPath, filesystem, log)
	reportGenerator := report.NewReportGenerator(stateManager, log, config.Chains)
	stateGenerator := state.NewStateGenerator(log, config.Chains)
	dataManager := data.NewManager(log, config.Chains)

	timeZone, _ := time.LoadLocation(config.Timezone)

	reporters := []reportersPkg.Reporter{
		pagerduty.NewPagerDutyReporter(config.PagerDutyConfig, log),
		telegram.NewTelegramReporter(config.TelegramConfig, mutesManager, stateGenerator, dataManager, log, version, timeZone),
	}

	reportDispatcher := report.NewDispatcher(log, mutesManager, reporters)

	return &App{
		Logger:           log,
		Config:           config,
		StateManager:     stateManager,
		ReportGenerator:  reportGenerator,
		StateGenerator:   stateGenerator,
		ReportDispatcher: reportDispatcher,
	}
}

func (a *App) Start() {
	a.StateManager.Load()
	if err := a.ReportDispatcher.Init(); err != nil {
		a.Logger.Panic().Err(err).Msg("Error initializing reporters")
	}

	c := cron.New()
	if _, err := c.AddFunc(a.Config.Interval, func() {
		a.Report()
	}); err != nil {
		a.Logger.Fatal().Err(err).Msg("Error processing cron pattern")
	}
	c.Start()
	a.Logger.Info().Str("interval", a.Config.Interval).Msg("Scheduled proposals reporting")

	select {}
}

func (a *App) Report() {
	newState := a.StateGenerator.GetState(a.StateManager.State)
	generatedReport := a.ReportGenerator.GenerateReport(a.StateManager.State, newState)
	a.ReportDispatcher.SendReport(generatedReport)
}
