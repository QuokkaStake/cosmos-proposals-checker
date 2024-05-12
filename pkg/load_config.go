package pkg

import (
	"main/pkg/fs"
	"main/pkg/types"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

func GetConfig(filesystem fs.FS, path string) (*types.Config, error) {
	configBytes, err := filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &types.Config{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	defaults.MustSet(configStruct)

	for _, chain := range configStruct.Chains {
		chain.Explorer = chain.GetExplorer()
	}

	return configStruct, nil
}
