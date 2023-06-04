package config

import (
	"fmt"
	"main/pkg/logger"
	"os"

	configTypes "main/pkg/config/types"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

type Config struct {
	PagerDutyConfig PagerDutyConfig       `toml:"pagerduty"`
	TelegramConfig  TelegramConfig        `toml:"telegram"`
	LogConfig       configTypes.LogConfig `toml:"log"`
	StatePath       string                `toml:"state-path"`
	MutesPath       string                `toml:"mutes-path"`
	Chains          configTypes.Chains    `toml:"chains"`
	Interval        string                `toml:"interval" default:"* * * * *"`
}

type PagerDutyConfig struct {
	PagerDutyURL string `toml:"url" default:"https://events.pagerduty.com"`
	APIKey       string `toml:"api-key"`
}

type TelegramConfig struct {
	TelegramChat  int64  `toml:"chat"`
	TelegramToken string `toml:"token"`
}

func (c *Config) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chains provided")
	}

	for index, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return fmt.Errorf("error in chain %d: %s", index, err)
		}
	}

	return nil
}

func GetConfig(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &Config{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	if err := defaults.Set(configStruct); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Error setting default config values")
	}

	for _, chain := range configStruct.Chains {
		if chain.MintscanPrefix != "" {
			chain.Explorer = &configTypes.Explorer{
				ProposalLinkPattern: fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", chain.MintscanPrefix),
				WalletLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/account/%%s", chain.MintscanPrefix),
			}
		}
	}

	return configStruct, nil
}
