package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateChainWithEmptyName(t *testing.T) {
	chain := Chain{
		Name: "",
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutEndpoints(t *testing.T) {
	chain := Chain{
		Name:         "chain",
		LCDEndpoints: []string{},
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutWallets(t *testing.T) {
	chain := Chain{
		Name:         "chain",
		LCDEndpoints: []string{"endpoint"},
		Wallets:      []string{},
	}

	err := chain.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithValidConfig(t *testing.T) {
	chain := Chain{
		Name:         "chain",
		LCDEndpoints: []string{"endpoint"},
		Wallets:      []string{"wallet"},
	}

	err := chain.Validate()
	assert.Equal(t, err, nil, "Error should not be presented!")
}

func TestChainGetNameWithoutPrettyName(t *testing.T) {
	chain := Chain{
		Name:       "chain",
		PrettyName: "",
	}

	err := chain.GetName()
	assert.Equal(t, err, "chain", "Chain name should match!")
}

func TestChainGetNameWithPrettyName(t *testing.T) {
	chain := Chain{
		Name:       "chain",
		PrettyName: "chain-pretty",
	}

	err := chain.GetName()
	assert.Equal(t, err, "chain-pretty", "Chain name should match!")
}

func TestValidateConfigNoChains(t *testing.T) {
	config := Config{
		Chains: []Chain{},
	}
	err := config.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidChain(t *testing.T) {
	config := Config{
		Chains: []Chain{
			{
				Name: "",
			},
		},
	}
	err := config.Validate()
	assert.NotEqual(t, err, nil, "Error should be presented!")
}

func TestValidateConfigValidChain(t *testing.T) {
	config := Config{
		Chains: []Chain{
			{
				Name:         "chain",
				LCDEndpoints: []string{"endpoint"},
				Wallets:      []string{"wallet"},
			},
		},
	}
	err := config.Validate()
	assert.Equal(t, err, nil, "Error should not be presented!")
}

func TestFindChainByNameIfPresent(t *testing.T) {
	chains := Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain2")
	assert.NotEqual(t, chain, nil, "Chain should be presented!")
}

func TestFindChainByNameIfNotPresent(t *testing.T) {
	chains := Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain3")
	assert.Nil(t, chain, "Chain should not be presented!")
}

func TestGetKeplrLink(t *testing.T) {
	chain := Chain{
		KeplrName: "chain",
	}

	link := chain.GetKeplrLink("proposal")
	assert.Equal(
		t,
		link,
		"https://wallet.keplr.app/#/chain/governance?detailId=proposal",
		"Chain Keplr link is wrong!",
	)
}
