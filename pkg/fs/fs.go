package fs

import (
	"io"
	"os"
)

type File interface {
	io.WriteCloser
}

type FS interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perms os.FileMode) error
	Create(path string) (File, error)
}
