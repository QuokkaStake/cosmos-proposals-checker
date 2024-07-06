package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOsFsRead(t *testing.T) {
	t.Parallel()

	fs := &OsFS{}
	file, err := fs.ReadFile("not-found.test")
	assert.Empty(t, file)
	require.Error(t, err)
}

func TestOsFsWrite(t *testing.T) {
	t.Parallel()

	fs := &OsFS{}
	err := fs.WriteFile("/etc/fstab", []byte{}, 0)
	require.Error(t, err)
}

func TestOsFsCreate(t *testing.T) {
	t.Parallel()

	fs := &OsFS{}
	_, err := fs.Create("/etc/fstab")
	require.Error(t, err)
}
