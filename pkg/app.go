package pkg

import (
	configPkg "main/pkg/config"
	"main/pkg/logger"
	"main/pkg/mutes_manager"
	reportPkg "main/pkg/report"
	reportersPkg "main/pkg/reporters"
	"main/pkg/reporters/pagerduty"
	"main/pkg/reporters/telegram"
	"main/pkg/state/generator"
	"main/pkg/state/manager"
	"time"
)

type App struct {
	ConfigPath string
}

func NewApp(configPath string) *App {
	return &App{ConfigPath: configPath}
}

func (a *App) Start() {
	config, err := configPkg.GetConfig(a.ConfigPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	log := logger.GetLogger(config.LogConfig)

	stateManager := manager.NewStateManager(config.StatePath, log)
	stateManager.Load()

	mutesManager := mutes_manager.NewMutesManager(config.MutesPath, log)
	mutesManager.Load()

	reportGenerator := reportPkg.NewReportGenerator(stateManager, log, config.Chains)
	stateGenerator := generator.NewStateGenerator(log, config.Chains)

	reporters := []reportersPkg.Reporter{
		pagerduty.NewPagerDutyReporter(config.PagerDutyConfig, log),
		telegram.NewTelegramReporter(config.TelegramConfig, mutesManager, stateGenerator, log),
	}

	for _, reporter := range reporters {
		reporter.Init()
		if reporter.Enabled() {
			log.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
	}

	for {
		newState := stateGenerator.GetState(stateManager.State)
		report := reportGenerator.GenerateReport(stateManager.State, newState)

		if report.Empty() {
			log.Debug().Msg("Empty report, not sending.")
			time.Sleep(time.Second * time.Duration(config.Interval))
			continue
		}

		log.Debug().Int("len", len(report.Entries)).Msg("Got non-empty report")

		for _, reporter := range reporters {
			if reporter.Enabled() {
				log.Debug().Str("name", reporter.Name()).Msg("Sending report...")
				if err := reporter.SendReport(report); err != nil {
					log.Error().Err(err).Str("name", reporter.Name()).Msg("Failed to send report")
				}
			}
		}

		time.Sleep(time.Second * time.Duration(config.Interval))
	}
}
