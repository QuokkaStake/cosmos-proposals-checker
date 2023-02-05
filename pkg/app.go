package pkg

import (
	configPkg "main/pkg/config"
	"main/pkg/logger"
	mutesManager "main/pkg/mutes_manager"
	reportPkg "main/pkg/report"
	reportersPkg "main/pkg/reporters"
	"main/pkg/reporters/pagerduty"
	"main/pkg/reporters/telegram"
	"main/pkg/state/generator"
	"main/pkg/state/manager"

	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type App struct {
	Logger          *zerolog.Logger
	Config          *configPkg.Config
	StateManager    *manager.StateManager
	MutesManager    *mutesManager.MutesManager
	ReportGenerator *reportPkg.ReportGenerator
	StateGenerator  *generator.StateGenerator
	Reporters       []reportersPkg.Reporter
}

func NewApp(configPath string) *App {
	config, err := configPkg.GetConfig(configPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	log := logger.GetLogger(config.LogConfig)

	stateManager := manager.NewStateManager(config.StatePath, log)
	mutesManager := mutesManager.NewMutesManager(config.MutesPath, log)
	reportGenerator := reportPkg.NewReportGenerator(stateManager, log, config.Chains)
	stateGenerator := generator.NewStateGenerator(log, config.Chains)

	return &App{
		Logger:          log,
		Config:          config,
		StateManager:    stateManager,
		MutesManager:    mutesManager,
		ReportGenerator: reportGenerator,
		StateGenerator:  stateGenerator,
		Reporters: []reportersPkg.Reporter{
			pagerduty.NewPagerDutyReporter(config.PagerDutyConfig, log),
			telegram.NewTelegramReporter(config.TelegramConfig, mutesManager, stateGenerator, log),
		},
	}
}

func (a *App) Start() {
	a.StateManager.Load()
	a.MutesManager.Load()

	for _, reporter := range a.Reporters {
		reporter.Init()
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
	report := a.ReportGenerator.GenerateReport(a.StateManager.State, newState)

	if report.Empty() {
		a.Logger.Debug().Msg("Empty report, not sending.")
		return
	}

	a.Logger.Debug().Int("len", len(report.Entries)).Msg("Got non-empty report")

	for _, reporter := range a.Reporters {
		if reporter.Enabled() {
			a.Logger.Debug().Str("name", reporter.Name()).Msg("Sending report...")
			if err := reporter.SendReport(report); err != nil {
				a.Logger.Error().Err(err).Str("name", reporter.Name()).Msg("Failed to send report")
			}
		}
	}
}
