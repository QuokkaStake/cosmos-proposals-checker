package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletAddressOrAliasWithoutAlias(t *testing.T) {
	t.Parallel()

	wallet := Wallet{Address: "test"}
	assert.Equal(t, "test", wallet.AddressOrAlias(), "Wrong value!")
}

func TestWalletAddressOrAliasWithAlias(t *testing.T) {
	t.Parallel()

	wallet := Wallet{Address: "test", Alias: "alias"}
	assert.Equal(t, "alias", wallet.AddressOrAlias(), "Wrong value!")
}

func TestGetLinksEmpty(t *testing.T) {
	t.Parallel()

	chain := Chain{}
	links := chain.GetExplorerProposalsLinks("test")

	assert.Empty(t, links, "Expected 0 links")
}

func TestGetLinksPresent(t *testing.T) {
	t.Parallel()

	chain := Chain{
		KeplrName: "chain",
		Explorer: &Explorer{
			ProposalLinkPattern: "example.com/proposal/%s",
		},
	}
	links := chain.GetExplorerProposalsLinks("test")

	assert.Len(t, links, 2, "Expected 2 links")
	assert.Equal(t, "Keplr", links[0].Name, "Expected Keplr link")
	assert.Equal(t, "https://wallet.keplr.app/chains/chain/proposals=test", links[0].Href, "Wrong Keplr link")
	assert.Equal(t, "Explorer", links[1].Name, "Expected Explorer link")
	assert.Equal(t, "example.com/proposal/test", links[1].Href, "Wrong explorer link")
}

func TestGetExplorerProposalLinkWithoutExplorer(t *testing.T) {
	t.Parallel()

	chain := Chain{}
	proposal := Proposal{ID: "ID", Title: "Title"}
	link := chain.GetProposalLink(proposal)

	assert.Equal(t, "Title", link.Name, "Wrong value!")
	assert.Equal(t, "", link.Href, "Wrong value!")
}

func TestGetExplorerProposalLinkWithExplorer(t *testing.T) {
	t.Parallel()

	chain := Chain{Explorer: &Explorer{ProposalLinkPattern: "example.com/%s"}}
	proposal := Proposal{ID: "ID", Title: "Title"}
	link := chain.GetProposalLink(proposal)

	assert.Equal(t, "Title", link.Name, "Wrong value!")
	assert.Equal(t, "example.com/ID", link.Href, "Wrong value!")
}

func TestGetWalletLinkWithoutExplorer(t *testing.T) {
	t.Parallel()

	chain := Chain{}
	wallet := &Wallet{Address: "wallet"}
	link := chain.GetWalletLink(wallet)

	assert.Equal(t, "wallet", link.Name, "Wrong value!")
	assert.Equal(t, "", link.Href, "Wrong value!")
}

func TestGetWalletLinkWithoutAlias(t *testing.T) {
	t.Parallel()

	chain := Chain{Explorer: &Explorer{WalletLinkPattern: "example.com/%s"}}
	wallet := &Wallet{Address: "wallet"}
	link := chain.GetWalletLink(wallet)

	assert.Equal(t, "wallet", link.Name, "Wrong value!")
	assert.Equal(t, "example.com/wallet", link.Href, "Wrong value!")
}

func TestGetWalletLinkWithAlias(t *testing.T) {
	t.Parallel()

	chain := Chain{Explorer: &Explorer{WalletLinkPattern: "example.com/%s"}}
	wallet := &Wallet{Address: "wallet", Alias: "alias"}
	link := chain.GetWalletLink(wallet)

	assert.Equal(t, "alias", link.Name, "Wrong value!")
	assert.Equal(t, "example.com/wallet", link.Href, "Wrong value!")
}

func TestSetExplorerMintscan(t *testing.T) {
	t.Parallel()

	chain := Chain{MintscanPrefix: "test"}
	explorer := chain.GetExplorer()

	assert.NotNil(t, explorer)
	assert.Equal(t, "https://mintscan.io/test/account/%s", explorer.WalletLinkPattern)
	assert.Equal(t, "https://mintscan.io/test/proposals/%s", explorer.ProposalLinkPattern)
}

func TestSetExplorerPing(t *testing.T) {
	t.Parallel()

	chain := Chain{PingPrefix: "test", PingHost: "https://example.com"}
	explorer := chain.GetExplorer()

	assert.NotNil(t, explorer)
	assert.Equal(t, "https://example.com/test/account/%s", explorer.WalletLinkPattern)
	assert.Equal(t, "https://example.com/test/gov/%s", explorer.ProposalLinkPattern)
}

func TestSetExplorerEmpty(t *testing.T) {
	t.Parallel()

	chain := Chain{}
	explorer := chain.GetExplorer()
	assert.Nil(t, explorer)
}

func TestChainDisplayWarningsEmptyExplorer(t *testing.T) {
	t.Parallel()

	chain := Chain{KeplrName: "test"}
	warnings := chain.DisplayWarnings()
	assert.Len(t, warnings, 1)
}

func TestChainDisplayWarningsEmptyKeplrName(t *testing.T) {
	t.Parallel()

	chain := Chain{Explorer: &Explorer{ProposalLinkPattern: "test", WalletLinkPattern: "test"}}
	warnings := chain.DisplayWarnings()
	assert.Len(t, warnings, 1)
}

func TestChainDisplayWarningsOk(t *testing.T) {
	t.Parallel()

	chain := Chain{
		Explorer:  &Explorer{ProposalLinkPattern: "test", WalletLinkPattern: "test"},
		KeplrName: "test",
	}
	warnings := chain.DisplayWarnings()
	assert.Empty(t, warnings)
}
