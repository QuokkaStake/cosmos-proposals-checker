package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Chain struct {
	Name                string   `toml:"name"`
	PrettyName          string   `toml:"pretty-name"`
	KeplrName           string   `toml:"keplr-name"`
	LCDEndpoints        []string `toml:"lcd-endpoints"`
	Wallets             []string `toml:"wallets"`
	MintscanPrefix      string   `toml:"mintscan-prefix"`
	ExplorerLinkPattern string   `toml:"explorer-link-pattern"`
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

func (c *Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c *Chain) GetKeplrLink(proposalID string) string {
	return fmt.Sprintf("https://wallet.keplr.app/#/%s/governance?detailId=%s", c.KeplrName, proposalID)
}

func (c *Chain) GetExplorerProposalsLinks(proposalID string) []ExplorerLink {
	if c.MintscanPrefix == "" {
		return []ExplorerLink{}
	}

	links := []ExplorerLink{
		{
			Name: "Mintscan",
			Link: fmt.Sprintf("https://mintscan.io/%s/proposals/%s", c.MintscanPrefix, proposalID),
		},
	}

	if c.ExplorerLinkPattern != "" {
		links = append(links, ExplorerLink{
			Name: "Explorer",
			Link: fmt.Sprintf(c.ExplorerLinkPattern, proposalID),
		})
	}

	return links
}

type Chains []Chain

func (c Chains) FindByName(name string) *Chain {
	for _, chain := range c {
		if chain.Name == name {
			return &chain
		}
	}

	return nil
}

type Config struct {
	PagerDutyConfig PagerDutyConfig `toml:"pagerduty"`
	TelegramConfig  TelegramConfig  `toml:"telegram"`
	LogConfig       LogConfig       `toml:"log"`
	StatePath       string          `toml:"state-path"`
	MutesPath       string          `toml:"mutes-path"`
	Chains          []Chain         `toml:"chains"`
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

	configStruct := Config{}
	if _, err = toml.Decode(configString, &configStruct); err != nil {
		return nil, err
	}

	defaults.SetDefaults(&configStruct)
	return &configStruct, nil
}
