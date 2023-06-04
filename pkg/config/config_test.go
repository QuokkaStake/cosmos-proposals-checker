package config

import (
	"testing"

	configTypes "main/pkg/config/types"

	"github.com/stretchr/testify/assert"
)

func TestValidateChainWithEmptyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name: "",
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutEndpoints(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:         "chain",
		LCDEndpoints: []string{},
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutWallets(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:         "chain",
		LCDEndpoints: []string{"endpoint"},
		Wallets:      []*configTypes.Wallet{},
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithValidConfig(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:          "chain",
		LCDEndpoints:  []string{"endpoint"},
		Wallets:       []*configTypes.Wallet{{Address: "wallet"}},
		ProposalsType: "v1",
	}

	err := chain.Validate()
	assert.Equal(t, err, nil, "Error should not be presented!")
}

func TestChainGetNameWithoutPrettyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:       "chain",
		PrettyName: "",
	}

	err := chain.GetName()
	assert.Equal(t, err, "chain", "Chain name should match!")
}

func TestChainGetNameWithPrettyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:       "chain",
		PrettyName: "chain-pretty",
	}

	err := chain.GetName()
	assert.Equal(t, err, "chain-pretty", "Chain name should match!")
}

func TestValidateConfigNoChains(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*configTypes.Chain{},
	}
	err := config.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidChain(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*configTypes.Chain{
			{
				Name: "",
			},
		},
	}
	err := config.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateConfigWrongProposalType(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*configTypes.Chain{
			{
				Name:          "chain",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*configTypes.Wallet{{Address: "wallet"}},
				ProposalsType: "test",
			},
		},
	}
	err := config.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateConfigValidChain(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*configTypes.Chain{
			{
				Name:          "chain",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*configTypes.Wallet{{Address: "wallet"}},
				ProposalsType: "v1",
			},
		},
	}
	err := config.Validate()
	assert.Equal(t, err, nil, "Error should not be presented!")
}

func TestFindChainByNameIfPresent(t *testing.T) {
	t.Parallel()

	chains := configTypes.Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain2")
	assert.NotEqual(t, chain, nil, "Chain should be presented!")
}

func TestFindChainByNameIfNotPresent(t *testing.T) {
	t.Parallel()

	chains := configTypes.Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain3")
	assert.Nil(t, chain, "Chain should not be presented!")
}

func TestGetLinksEmpty(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{}
	links := chain.GetExplorerProposalsLinks("test")

	assert.Equal(t, len(links), 0, "Expected 0 links")
}

func TestGetLinksPresent(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		KeplrName: "chain",
		Explorer: &configTypes.Explorer{
			ProposalLinkPattern: "example.com/proposal/%s",
		},
	}
	links := chain.GetExplorerProposalsLinks("test")

	assert.Equal(t, len(links), 2, "Expected 2 links")
	assert.Equal(t, links[0].Name, "Keplr", "Expected Keplr link")
	assert.Equal(t, links[0].Href, "https://wallet.keplr.app/#/chain/governance?detailId=test", "Wrong Keplr link")
	assert.Equal(t, links[1].Name, "Explorer", "Expected Explorer link")
	assert.Equal(t, links[1].Href, "example.com/proposal/test", "Wrong explorer link")
}
