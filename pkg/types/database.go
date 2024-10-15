package types

import (
	"errors"
)

type DatabaseConfig struct {
	Path string `toml:"path"`
}

func (c *DatabaseConfig) Validate() error {
	if c.Path == "" {
		return errors.New("database path not specified")
	}

	return nil
}
