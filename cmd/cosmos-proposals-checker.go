package main

import (
	"main/pkg"
	"main/pkg/logger"

	"github.com/spf13/cobra"
)

func Execute(configPath string) {
	app := pkg.NewApp(configPath)
	app.Start()
}

var (
	version = "unknown"
)

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:     "cosmos-proposals-checker",
		Long:    "Checks the specific wallets on different chains for proposal votes.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			Execute(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
