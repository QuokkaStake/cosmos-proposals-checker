package types

import (
	"fmt"
	"main/pkg/types"
)

type Explorer struct {
	ProposalLinkPattern string `toml:"proposal-link-pattern"`
	WalletLinkPattern   string `toml:"wallet-link-pattern"`
}

type Wallet struct {
	Address string `toml:"address"`
	Alias   string `toml:"alias"`
}

func (w *Wallet) AddressOrAlias() string {
	if w.Alias != "" {
		return w.Alias
	}

	return w.Address
}

type Chain struct {
	Name           string    `toml:"name"`
	PrettyName     string    `toml:"pretty-name"`
	KeplrName      string    `toml:"keplr-name"`
	LCDEndpoints   []string  `toml:"lcd-endpoints"`
	Wallets        []*Wallet `toml:"wallets"`
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

	for index, wallet := range c.Wallets {
		if wallet.Address == "" {
			return fmt.Errorf("wallet #%d: address is empty", index)
		}
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

func (c Chain) GetProposalLink(proposal types.Proposal) types.Link {
	if c.Explorer == nil || c.Explorer.ProposalLinkPattern == "" {
		return types.Link{Name: proposal.Content.Title}
	}

	return types.Link{
		Name: proposal.Content.Title,
		Href: fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposal.ProposalID),
	}
}

func (c Chain) GetWalletLink(wallet *Wallet) types.Link {
	if c.Explorer == nil || c.Explorer.WalletLinkPattern == "" {
		return types.Link{Name: wallet.Address}
	}

	link := types.Link{
		Name: wallet.Address,
		Href: fmt.Sprintf(c.Explorer.WalletLinkPattern, wallet.Address),
	}

	if wallet.Alias != "" {
		link.Name = wallet.Alias
	}

	return link
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
