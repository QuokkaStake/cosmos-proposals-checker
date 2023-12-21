package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeLinkWithoutHref(t *testing.T) {
	t.Parallel()

	link := Link{
		Name: "link",
	}

	assert.Equal(t, "link", link.Serialize(), "Wrong value!")
}

func TestSerializeLinkWithHref(t *testing.T) {
	t.Parallel()

	link := Link{
		Name: "link",
		Href: "example.com",
	}

	assert.Equal(t, "<a href='example.com'>link</a>", link.Serialize(), "Wrong value!")
}
