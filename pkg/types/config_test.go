package types

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestValidateChainWithEmptyName(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name: "",
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutEndpoints(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name:         "chain",
		Type:         "cosmos",
		LCDEndpoints: []string{},
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithoutWallets(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name:         "chain",
		Type:         "cosmos",
		LCDEndpoints: []string{"endpoint"},
		Wallets:      []*Wallet{},
	}

	err := chain.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateChainWithValidConfig(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name:          "chain",
		LCDEndpoints:  []string{"endpoint"},
		Wallets:       []*Wallet{{Address: "wallet"}},
		ProposalsType: "v1",
		Type:          "cosmos",
	}

	err := chain.Validate()
	require.NoError(t, err, "Error should not be presented!")
}

func TestChainGetNameWithoutPrettyName(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name:       "chain",
		PrettyName: "",
	}

	name := chain.GetName()
	assert.Equal(t, "chain", name, "Chain name should match!")
}

func TestChainGetNameWithPrettyName(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Name:       "chain",
		PrettyName: "chain-pretty",
	}

	err := chain.GetName()
	assert.Equal(t, "chain-pretty", err, "Chain name should match!")
}

func TestValidateConfigNoDatabase(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*Chain{},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigNoChains(t *testing.T) {
	t.Parallel()

	config := Config{
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains:         []*Chain{},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidChain(t *testing.T) {
	t.Parallel()

	config := Config{
		Chains: []*Chain{
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
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains: []*Chain{
			{
				Name:          "chain",
				Type:          "cosmos",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*Wallet{{Address: "wallet"}},
				ProposalsType: "test",
			},
		},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidTimezone(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "test",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains: []*Chain{
			{
				Name:          "chain",
				Type:          "cosmos",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*Wallet{{Address: "wallet"}},
				ProposalsType: "v1",
			},
		},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidWallet(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains: []*Chain{
			{
				Name:          "chain",
				Type:          "cosmos",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*Wallet{{Address: ""}},
				ProposalsType: "v1",
			},
		},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigInvalidType(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains: []*Chain{
			{
				Name:          "chain",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*Wallet{{Address: "wallet"}},
				ProposalsType: "v1",
				Type:          "invalid",
			},
		},
	}
	err := config.Validate()
	require.Error(t, err, nil, "Error should be presented!")
}

func TestValidateConfigValidChain(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		Chains: []*Chain{
			{
				Name:          "chain",
				LCDEndpoints:  []string{"endpoint"},
				Wallets:       []*Wallet{{Address: "wallet"}},
				ProposalsType: "v1",
				Type:          "cosmos",
			},
		},
	}
	err := config.Validate()
	require.NoError(t, err, "Error should not be presented!")
}

func TestConfigDisplayWarningInvalidChain(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		StatePath:      "test",
		MutesPath:      "test",
		Chains: []*Chain{
			{
				KeplrName: "test",
			},
		},
	}
	warnings := config.DisplayWarnings()
	assert.Len(t, warnings, 1)
}

func TestConfigDisplayWarningNoStatePath(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		MutesPath:      "test",
		Chains: []*Chain{
			{
				KeplrName: "test",
				Explorer:  &Explorer{WalletLinkPattern: "test", ProposalLinkPattern: "test"},
			},
		},
	}
	warnings := config.DisplayWarnings()
	assert.Len(t, warnings, 1)
}

func TestConfigDisplayWarningNoMutesPath(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		StatePath:      "test",
		Chains: []*Chain{
			{
				KeplrName: "test",
				Explorer:  &Explorer{WalletLinkPattern: "test", ProposalLinkPattern: "test"},
			},
		},
	}
	warnings := config.DisplayWarnings()
	assert.Len(t, warnings, 1)
}

func TestConfigDisplayWarningOk(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone:       "Europe/Moscow",
		DatabaseConfig: DatabaseConfig{Path: "database.sqlite"},
		StatePath:      "test",
		MutesPath:      "test",
		Chains: []*Chain{
			{
				KeplrName: "test",
				Explorer:  &Explorer{WalletLinkPattern: "test", ProposalLinkPattern: "test"},
			},
		},
	}
	warnings := config.DisplayWarnings()
	assert.Empty(t, warnings)
}

func TestConfigLogWarnings(t *testing.T) {
	t.Parallel()

	config := Config{
		Timezone: "Europe/Moscow",
		Chains: []*Chain{
			{},
		},
	}
	warnings := config.DisplayWarnings()
	logger := zerolog.Nop()
	config.LogWarnings(&logger, warnings)
}
