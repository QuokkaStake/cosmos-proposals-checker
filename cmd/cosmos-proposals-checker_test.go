package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled
func TestValidateConfigNoConfigProvided(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	os.Args = []string{"cmd", "validate-config"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestValidateConfigFailedToLoad(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-not-found.toml"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestValidateConfigInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-invalid.toml"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestValidateConfigWithWarnings(t *testing.T) {
	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-with-warnings.toml"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestValidateConfigValid(t *testing.T) {
	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-valid.toml"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestStartNoConfigProvided(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	os.Args = []string{"cmd"}
	main()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestStartConfigProvided(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	os.Args = []string{"cmd", "--config", "../assets/config-invalid.toml"}
	main()
	assert.True(t, true)
}
