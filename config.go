package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Chain struct {
	Name         string   `toml:"name"`
	LCDEndpoints []string `toml:"lcd-endpoints"`
	Wallets      []string `toml:"wallets"`
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if len(c.LCDEndpoints) == 0 {
		return fmt.Errorf("no LCD endpoints provided")
	}

	if len(c.Wallets) == 0 {
		return fmt.Errorf("no wallets provided")
	}

	return nil
}

type Config struct {
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	LogConfig       LogConfig       `toml:"log"`
	StatePath       string          `toml:"state-path"`
	Chains          []Chain         `toml:"chains"`
}

type PagerDutyConfig struct {
	PagerDutyURL string `toml:"url"`
	APIKey       string `toml:"api-key"`
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

	configStruct := Config{}
	if _, err = toml.Decode(configString, &configStruct); err != nil {
		return nil, err
	}

	defaults.SetDefaults(&configStruct)
	return &configStruct, nil
}
