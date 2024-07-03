package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileWrite(t *testing.T) {
	t.Parallel()

	file := TestFile{WithFileWriteError: true}
	_, err := file.Write([]byte{})
	require.Error(t, err)

	file2 := TestFile{}
	_, err2 := file2.Write([]byte{})
	require.NoError(t, err2)
}

func TestFileClose(t *testing.T) {
	t.Parallel()

	file := TestFile{WithFileCloseError: true}
	err := file.Close()
	require.Error(t, err)

	file2 := TestFile{}
	err2 := file2.Close()
	require.NoError(t, err2)
}

func TestFsRead(t *testing.T) {
	t.Parallel()

	fs := &TestFS{}
	file, err := fs.ReadFile("lcd-error.json")
	assert.NotEmpty(t, file)
	require.NoError(t, err)
}

func TestFsWrite(t *testing.T) {
	t.Parallel()

	fs := &TestFS{WithWriteError: true}
	err := fs.WriteFile("lcd-error.json", []byte{}, 0)
	require.Error(t, err)

	fs2 := &TestFS{}
	err2 := fs2.WriteFile("lcd-error.json", []byte{}, 0)
	require.NoError(t, err2)
}

func TestFsCreate(t *testing.T) {
	t.Parallel()

	fs := &TestFS{WithCreateError: true}
	_, err := fs.Create("lcd-error.json")
	require.Error(t, err)

	fs2 := &TestFS{}
	_, err2 := fs2.Create("lcd-error.json")
	require.NoError(t, err2)
}
