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

func Execute(configPath string) {
	filesystem := &OsFS{}
	app := pkg.NewApp(configPath, filesystem, version)
	app.Start()
}

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
