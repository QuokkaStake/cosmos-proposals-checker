package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
