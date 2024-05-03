package fs

import (
	"errors"
	"main/assets"
	"os"
)

type TestFS struct {
	WithWriteError  bool
	WithCreateError bool

	WithFileWriteError bool
	WithFileCloseError bool
}

type TestFile struct {
	WithFileWriteError bool
	WithFileCloseError bool
}

func (f *TestFile) Write(p []byte) (int, error) {
	if f.WithFileWriteError {
		return 0, errors.New("stub error")
	}

	return 0, nil
}

func (f *TestFile) Close() error {
	if f.WithFileCloseError {
		return errors.New("stub error")
	}

	return nil
}

func (fs *TestFS) ReadFile(name string) ([]byte, error) {
	return assets.EmbedFS.ReadFile(name)
}

func (fs *TestFS) WriteFile(name string, data []byte, perms os.FileMode) error {
	if fs.WithWriteError {
		return errors.New("stub error")
	}

	return nil
}

func (fs *TestFS) Create(path string) (File, error) {
	if fs.WithCreateError {
		return nil, errors.New("stub error")
	}

	return &TestFile{
		WithFileWriteError: fs.WithFileWriteError,
		WithFileCloseError: fs.WithFileCloseError,
	}, nil
}
