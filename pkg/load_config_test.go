package pkg

import (
	"main/pkg/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigNotExistingFile(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	config, err := GetConfig(filesystem, "notexisting.toml")

	assert.Nil(t, config)
	require.Error(t, err)
}

func TestLoadConfigInvalidToml(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	config, err := GetConfig(filesystem, "invalid.toml")

	assert.Nil(t, config)
	require.Error(t, err)
}

func TestLoadConfigValidToml(t *testing.T) {
	t.Parallel()

	filesystem := &fs.TestFS{}

	config, err := GetConfig(filesystem, "config-valid.toml")

	require.NoError(t, err)
	assert.NotNil(t, config)
	require.Len(t, config.Chains, 1)

	firstChain := config.Chains[0]
	require.NotNil(t, firstChain.Explorer)
	require.Equal(t, "https://mintscan.io/bitsong/proposals/%s", firstChain.Explorer.ProposalLinkPattern)
	require.Equal(t, "https://mintscan.io/bitsong/account/%s", firstChain.Explorer.WalletLinkPattern)
}
