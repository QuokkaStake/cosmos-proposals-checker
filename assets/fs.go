package assets

import (
	"embed"
)

//go:embed *
var EmbedFS embed.FS

func GetBytesOrPanic(path string) []byte {
	bytes, err := EmbedFS.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return bytes
}
