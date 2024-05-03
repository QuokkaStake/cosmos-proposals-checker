package pkg

import (
	"fmt"
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
		if chain.MintscanPrefix != "" {
			chain.Explorer = &types.Explorer{
				ProposalLinkPattern: fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", chain.MintscanPrefix),
				WalletLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/account/%%s", chain.MintscanPrefix),
			}
		}
	}

	return configStruct, nil
}
