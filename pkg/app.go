package pkg

import (
	"main/pkg/data"
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
	Logger          *zerolog.Logger
	Config          *types.Config
	StateManager    *state.Manager
	MutesManager    *mutes.Manager
	ReportGenerator *report.Generator
	StateGenerator  *state.Generator
	Reporters       []reportersPkg.Reporter
}

func NewApp(configPath string, version string) *App {
	config, err := GetConfig(configPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	log := logger.GetLogger(config.LogConfig)

	stateManager := state.NewStateManager(config.StatePath, log)
	mutesManager := mutes.NewMutesManager(config.MutesPath, log)
	reportGenerator := report.NewReportGenerator(stateManager, log, config.Chains)
	stateGenerator := state.NewStateGenerator(log, config.Chains)
	dataManager := data.NewManager(log, config.Chains)

	timeZone, _ := time.LoadLocation(config.Timezone)

	return &App{
		Logger:          log,
		Config:          config,
		StateManager:    stateManager,
		MutesManager:    mutesManager,
		ReportGenerator: reportGenerator,
		StateGenerator:  stateGenerator,
		Reporters: []reportersPkg.Reporter{
			pagerduty.NewPagerDutyReporter(config.PagerDutyConfig, log),
			telegram.NewTelegramReporter(config.TelegramConfig, mutesManager, stateGenerator, dataManager, log, version, timeZone),
		},
	}
}

func (a *App) Start() {
	a.StateManager.Load()
	a.MutesManager.Load()

	for _, reporter := range a.Reporters {
		if err := reporter.Init(); err != nil {
			a.Logger.Fatal().Err(err).Str("name", reporter.Name()).Msg("Error initializing reporter")
		}
		if reporter.Enabled() {
			a.Logger.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
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

	if generatedReport.Empty() {
		a.Logger.Debug().Msg("Empty report, not sending.")
		return
	}

	a.Logger.Debug().Int("len", len(generatedReport.Entries)).Msg("Got non-empty report")

	for _, reporter := range a.Reporters {
		if reporter.Enabled() {
			a.Logger.Debug().Str("name", reporter.Name()).Msg("Sending report...")
			if err := reporter.SendReport(generatedReport); err != nil {
				a.Logger.Error().Err(err).Str("name", reporter.Name()).Msg("Failed to send report")
			}
		}
	}
}
