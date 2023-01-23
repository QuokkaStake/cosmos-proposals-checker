package types

import (
	"fmt"
	"html/template"
	"main/pkg/types"
)

type Explorer struct {
	ProposalLinkPattern string `toml:"proposal-link-pattern"`
	WalletLinkPattern   string `toml:"wallet-link-pattern"`
}

type Chain struct {
	Name           string    `toml:"name"`
	PrettyName     string    `toml:"pretty-name"`
	KeplrName      string    `toml:"keplr-name"`
	LCDEndpoints   []string  `toml:"lcd-endpoints"`
	Wallets        []string  `toml:"wallets"`
	MintscanPrefix string    `toml:"mintscan-prefix"`
	Explorer       *Explorer `toml:"explorer"`
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

func (c Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c Chain) GetExplorerProposalsLinks(proposalID string) []types.Link {
	links := []types.Link{}

	if c.KeplrName != "" {
		links = append(links, types.Link{
			Name: "Keplr",
			Href: fmt.Sprintf("https://wallet.keplr.app/#/%s/governance?detailId=%s", c.KeplrName, proposalID),
		})
	}

	if c.Explorer != nil && c.Explorer.ProposalLinkPattern != "" {
		links = append(links, types.Link{
			Name: "Explorer",
			Href: fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposalID),
		})
	}

	return links
}

func (c Chain) GetWalletLink(wallet string) template.HTML {
	if c.Explorer == nil || c.Explorer.WalletLinkPattern == "" {
		return template.HTML(wallet)
	}

	link := fmt.Sprintf(c.Explorer.WalletLinkPattern, wallet)
	return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link, wallet))
}

type Chains []*Chain

func (c Chains) FindByName(name string) *Chain {
	for _, chain := range c {
		if chain.Name == name {
			return chain
		}
	}

	return nil
}
