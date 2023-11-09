package types

import "fmt"

type Config struct {
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	TelegramConfig  TelegramConfig  `toml:"telegram"`
	LogConfig       LogConfig       `toml:"log"`
	StatePath       string          `toml:"state-path"`
	MutesPath       string          `toml:"mutes-path"`
	Chains          Chains          `toml:"chains"`
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

	return nil
}
