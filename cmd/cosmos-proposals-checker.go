package main

import (
	"main/pkg"
	"main/pkg/fs"
	"main/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

type OsFS struct {
}

func (fs *OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OsFS) WriteFile(name string, data []byte, perms os.FileMode) error {
	return os.WriteFile(name, data, perms)
}

func (fs *OsFS) Create(path string) (fs.File, error) {
	return os.Create(path)
}

func ExecuteMain(configPath string) {
	filesystem := &OsFS{}
	app := pkg.NewApp(configPath, filesystem, version)
	app.Start()
}

func ExecuteValidateConfig(configPath string) {
	filesystem := &OsFS{}

	config, err := pkg.GetConfig(filesystem, configPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config!")
	}

	if err := config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Config is invalid!")
	}

	if warnings := config.DisplayWarnings(); len(warnings) > 0 {
		config.LogWarnings(logger.GetDefaultLogger(), warnings)
	} else {
		logger.GetDefaultLogger().Info().Msg("Provided config is valid.")
	}
}

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:     "cosmos-proposals-checker --config [config path]",
		Long:    "Checks the specific wallets on different chains for proposal votes.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteMain(ConfigPath)
		},
	}

	validateConfigCmd := &cobra.Command{
		Use:     "validate-config --config [config path]",
		Long:    "Checks the specific wallets on different chains for proposal votes.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteValidateConfig(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	validateConfigCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := validateConfigCmd.MarkPersistentFlagRequired("config"); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	rootCmd.AddCommand(validateConfigCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
