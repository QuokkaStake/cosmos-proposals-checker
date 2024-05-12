package types

type Explorer struct {
	ProposalLinkPattern string `toml:"proposal-link-pattern"`
	WalletLinkPattern   string `toml:"wallet-link-pattern"`
}
