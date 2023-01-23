package config

import (
	"fmt"
	"main/pkg/config/types"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	TelegramConfig  TelegramConfig  `toml:"telegram"`
	LogConfig       LogConfig       `toml:"log"`
	StatePath       string          `toml:"state-path"`
	MutesPath       string          `toml:"mutes-path"`
	Chains          types.Chains    `toml:"chains"`
	Interval        int64           `toml:"interval" default:"3600"`
}

type PagerDutyConfig struct {
	PagerDutyURL string `toml:"url" default:"https://events.pagerduty.com"`
	APIKey       string `toml:"api-key"`
}

type TelegramConfig struct {
	TelegramChat  int64  `toml:"chat"`
	TelegramToken string `toml:"token"`
}

type LogConfig struct {
	LogLevel   string `toml:"level" default:"info"`
	JSONOutput bool   `toml:"json" default:"false"`
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

	defaults.SetDefaults(configStruct)

	for _, chain := range configStruct.Chains {
		if chain.MintscanPrefix != "" {
			chain.Explorer = &types.Explorer{
				ProposalLinkPattern: fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", chain.MintscanPrefix),
				WalletLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/account/%%s", chain.MintscanPrefix),
			}
		}
	}

	return configStruct, nil
}
