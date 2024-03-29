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

func TestFindChainByNameIfPresent(t *testing.T) {
	t.Parallel()

	chains := Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain2")
	assert.NotNil(t, chain, "Chain should be presented!")
}

func TestFindChainByNameIfNotPresent(t *testing.T) {
	t.Parallel()

	chains := Chains{
		{Name: "chain1"},
		{Name: "chain2"},
	}

	chain := chains.FindByName("chain3")
	assert.Nil(t, chain, "Chain should not be presented!")
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
	assert.Equal(t, "https://wallet.keplr.app/#/chain/governance?detailId=test", links[0].Href, "Wrong Keplr link")
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
