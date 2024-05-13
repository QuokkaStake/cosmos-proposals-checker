package types

type Explorer struct {
	ProposalLinkPattern string `toml:"proposal-link-pattern"`
	WalletLinkPattern   string `toml:"wallet-link-pattern"`
}

func (e *Explorer) DisplayWarnings(chainName string) []Warning {
	warnings := make([]Warning, 0)

	if e.WalletLinkPattern == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": chainName},
			Message: "wallet-link-pattern for explorer is not set, cannot generate wallet links",
		})
	}

	if e.ProposalLinkPattern == "" {
		warnings = append(warnings, Warning{
			Labels:  map[string]string{"chain": chainName},
			Message: "proposal-link-pattern for explorer is not set, cannot generate proposal links",
		})
	}

	return warnings
}
