package fs

import (
	"main/assets"
	"os"
)

type TestFS struct{}

type TestFile struct {
}

func (f *TestFile) Write(p []byte) (int, error) {
	return 0, nil
}

func (f *TestFile) Close() error {
	return nil
}

func (fs *TestFS) ReadFile(name string) ([]byte, error) {
	return assets.EmbedFS.ReadFile(name)
}

func (fs *TestFS) WriteFile(name string, data []byte, perms os.FileMode) error {
	return nil
}

func (fs *TestFS) Create(path string) (File, error) {
	return &TestFile{}, nil // go
}
