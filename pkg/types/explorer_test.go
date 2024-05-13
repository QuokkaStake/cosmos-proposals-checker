package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExplorerDisplayWarningsNoProposalsLinkPattern(t *testing.T) {
	t.Parallel()

	explorer := &Explorer{WalletLinkPattern: "test"}
	warnings := explorer.DisplayWarnings("test")
	assert.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsNoWalletLinkPattern(t *testing.T) {
	t.Parallel()

	explorer := &Explorer{ProposalLinkPattern: "test"}
	warnings := explorer.DisplayWarnings("test")
	assert.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsOk(t *testing.T) {
	t.Parallel()

	explorer := &Explorer{ProposalLinkPattern: "test", WalletLinkPattern: "test"}
	warnings := explorer.DisplayWarnings("test")
	assert.Empty(t, warnings)
}
