package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	configTypes "main/pkg/config/types"

	"github.com/stretchr/testify/assert"
)

func TestValidateChainWithEmptyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name: "",
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutEndpoints(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:         "chain",
		LCDEndpoints: []string{},
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutWallets(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:         "chain",
		LCDEndpoints: []string{"endpoint"},
		Wallets:      []*configTypes.Wallet{},
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
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
	require.NoError(t, err, "Error should not be presented!")
}

func TestChainGetNameWithoutPrettyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:       "chain",
		PrettyName: "",
	}

	name := chain.GetName()
	assert.Equal(t, "chain", name, "Chain name should match!")
}

func TestChainGetNameWithPrettyName(t *testing.T) {
	t.Parallel()

	chain := configTypes.Chain{
		Name:       "chain",
		PrettyName: "chain-pretty",
	}

	err := chain.GetName()
	assert.Equal(t, "chain-pretty", err, "Chain name should match!")
}

func TestValidateConfigNoChains(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*configTypes.Chain{},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
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
	require.Error(t, err, nil, "Error should be presented!")
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
	require.Error(t, err, nil, "Error should be presented!")
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
	require.NoError(t, err, "Error should not be presented!")
}

func TestFindChainByNameIfPresent(t *testing.T) {
	t.Parallel()

	chains := configTypes.Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain2")
	assert.NotNil(t, chain, "Chain should be presented!")
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

	assert.Empty(t, links, "Expected 0 links")
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

	assert.Len(t, links, 2, "Expected 2 links")
	assert.Equal(t, "Keplr", links[0].Name, "Expected Keplr link")
	assert.Equal(t, "https://wallet.keplr.app/#/chain/governance?detailId=test", links[0].Href, "Wrong Keplr link")
	assert.Equal(t, "Explorer", links[1].Name, "Expected Explorer link")
	assert.Equal(t, "example.com/proposal/test", links[1].Href, "Wrong explorer link")
}
