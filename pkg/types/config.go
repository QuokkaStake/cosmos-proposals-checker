package types

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	TelegramConfig  TelegramConfig  `toml:"telegram"`
	LogConfig       LogConfig       `toml:"log"`
	StatePath       string          `toml:"state-path"`
	MutesPath       string          `toml:"mutes-path"`
	Chains          Chains          `toml:"chains"`
	Timezone        string          `toml:"timezone"`
	Interval        string          `default:"* * * * *" toml:"interval"`
}

type PagerDutyConfig struct {
	PagerDutyURL string `default:"https://events.pagerduty.com" toml:"url"`
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

	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return fmt.Errorf("error parsing timezone: %s", err)
	}

	return nil
}

func (c *Config) DisplayWarnings() []Warning {
	warnings := make([]Warning, 0)

	for _, chain := range c.Chains {
		warnings = append(warnings, chain.DisplayWarnings()...)
	}

	if c.MutesPath == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{},
			Message: "mutes-path is not set, cannot persist proposals mutes on disk.",
		})
	}

	if c.StatePath == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{},
			Message: "state-path is not set, cannot persist proposals state on disk.",
		})
	}

	return warnings
}

func (c *Config) LogWarnings(logger *zerolog.Logger, warnings []Warning) {
	for _, warning := range warnings {
		entry := logger.Warn()

		for key, label := range warning.Labels {
			entry = entry.Str(key, label)
		}

		entry.Msg(warning.Message)
	}
}
