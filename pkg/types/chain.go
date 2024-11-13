package types

import (
	"fmt"
	"main/pkg/utils"
)

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
	ProposalsType  string    `default:"v1beta1"          toml:"proposals-type"`
	Wallets        []*Wallet `toml:"wallets"`
	MintscanPrefix string    `toml:"mintscan-prefix"`
	PingPrefix     string    `toml:"ping-prefix"`
	PingHost       string    `default:"https://ping.pub" toml:"ping-host"`
	Explorer       *Explorer `toml:"explorer"`

	Type                 string `default:"cosmos"                                                             toml:"type"`
	NeutronSmartContract string `default:"neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh" toml:"neutron-smart-contract"`
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if !utils.Contains([]string{"cosmos", "neutron"}, c.Type) {
		return fmt.Errorf("expected type to be one of 'cosmos', 'neutron', but got '%s'", c.Type)
	}

	if len(c.LCDEndpoints) == 0 {
		return fmt.Errorf("no LCD endpoints provided")
	}

	if len(c.Wallets) == 0 {
		return fmt.Errorf("no wallets provided")
	}

	if !utils.Contains([]string{"v1beta1", "v1"}, c.ProposalsType) {
		return fmt.Errorf("wrong proposals type: expected one of 'v1beta1', 'v1', but got %s", c.ProposalsType)
	}

	for index, wallet := range c.Wallets {
		if wallet.Address == "" {
			return fmt.Errorf("wallet #%d: address is empty", index)
		}
	}

	return nil
}

func (c *Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c *Chain) GetExplorerProposalsLinks(proposalID string) []Link {
	links := []Link{}

	if c.KeplrName != "" {
		links = append(links, Link{
			Name: "Keplr",
			Href: fmt.Sprintf("https://wallet.keplr.app/chains/%s/proposals/%s", c.KeplrName, proposalID),
		})
	}

	if c.Explorer != nil && c.Explorer.ProposalLinkPattern != "" {
		links = append(links, Link{
			Name: "Explorer",
			Href: fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposalID),
		})
	}

	return links
}

func (c *Chain) GetProposalLink(proposal Proposal) Link {
	if c.Explorer == nil || c.Explorer.ProposalLinkPattern == "" {
		return Link{Name: proposal.Title}
	}

	return Link{
		Name: proposal.Title,
		Href: fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposal.ID),
	}
}

func (c *Chain) GetWalletLink(wallet *Wallet) Link {
	if c.Explorer == nil || c.Explorer.WalletLinkPattern == "" {
		return Link{Name: wallet.AddressOrAlias()}
	}

	link := Link{
		Name: wallet.AddressOrAlias(),
		Href: fmt.Sprintf(c.Explorer.WalletLinkPattern, wallet.Address),
	}

	return link
}

func (c *Chain) GetExplorer() *Explorer {
	if c.MintscanPrefix != "" {
		return &Explorer{
			ProposalLinkPattern: fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", c.MintscanPrefix),
			WalletLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/account/%%s", c.MintscanPrefix),
		}
	}

	if c.PingPrefix != "" {
		return &Explorer{
			ProposalLinkPattern: fmt.Sprintf("%s/%s/gov/%%s", c.PingHost, c.PingPrefix),
			WalletLinkPattern:   fmt.Sprintf("%s/%s/account/%%s", c.PingHost, c.PingPrefix),
		}
	}

	return c.Explorer
}

func (c *Chain) DisplayWarnings() []Warning {
	warnings := make([]Warning, 0)

	if c.Explorer == nil {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": c.Name},
			Message: "explorer is not set, cannot generate links",
		})
	} else {
		warnings = append(warnings, c.Explorer.DisplayWarnings(c.Name)...)
	}

	if c.KeplrName == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": c.Name},
			Message: "keplr-name is not set, cannot generate Keplr link to proposal",
		})
	}

	return warnings
}
