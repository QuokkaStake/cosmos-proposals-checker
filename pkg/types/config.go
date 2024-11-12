package types

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	DatabaseConfig  DatabaseConfig  `toml:"database"`
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	TelegramConfig  TelegramConfig  `toml:"telegram"`
	DiscordConfig   DiscordConfig   `toml:"discord"`
	LogConfig       LogConfig       `toml:"log"`
	TracingConfig   TracingConfig   `toml:"tracing"`
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

type DiscordConfig struct {
	Guild   string `toml:"guild"`
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

func (c *Config) Validate() error {
	if err := c.DatabaseConfig.Validate(); err != nil {
		return fmt.Errorf("invalid database config: %s", err)
	}

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
