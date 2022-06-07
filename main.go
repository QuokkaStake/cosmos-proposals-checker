package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
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
	reportGenerator := NewReportGenerator(stateManager, log, config.Chains)

	for {
		_ = reportGenerator.GenerateReport()
		time.Sleep(time.Second * 30)
	}
}

func main() {
	var ConfigPath string

	var rootCmd = &cobra.Command{
		Use:  "cosmos-proposals-checker",
		Long: "Checks the specific wallets on different chains for proposal votes.",
		Run: func(cmd *cobra.Command, args []string) {
			Execute(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	rootCmd.MarkPersistentFlagRequired("config")

	if err := rootCmd.Execute(); err != nil {
		log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		log.Fatal().Err(err).Msg("Could not start application")
	}
}
