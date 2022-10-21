package main

import (
	"time"

	"github.com/spf13/cobra"
)

const PaginationLimit = 1000

func Execute(configPath string) {
	config, err := GetConfig(configPath)
	if err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	log := GetLogger(config.LogConfig)

	stateManager := NewStateManager(config.StatePath, log)
	stateManager.Load()

	mutesManager := NewMutesManager(config.MutesPath, log)
	mutesManager.Load()

	reportGenerator := NewReportGenerator(stateManager, log, config.Chains)
	stateGenerator := NewStateGenerator(log, config.Chains)

	reporters := []Reporter{
		NewPagerDutyReporter(config.PagerDutyConfig, log),
		NewTelegramReporter(config.TelegramConfig, mutesManager, stateGenerator, log),
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

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:  "cosmos-proposals-checker",
		Long: "Checks the specific wallets on different chains for proposal votes.",
		Run: func(cmd *cobra.Command, args []string) {
			Execute(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
