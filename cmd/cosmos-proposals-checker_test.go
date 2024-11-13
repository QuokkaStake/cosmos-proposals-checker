package main

import (
	"os"
	"testing"

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
}

//nolint:paralleltest // disabled
func TestValidateConfigWithWarnings(t *testing.T) {
	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-with-warnings.toml"}
	main()

	require.NotEmpty(t, os.Args)
}

//nolint:paralleltest // disabled
func TestValidateConfigValid(t *testing.T) {
	os.Args = []string{"cmd", "validate-config", "--config", "../assets/config-valid.toml"}
	main()

	require.NotEmpty(t, os.Args)
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
}
