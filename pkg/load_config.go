package pkg

import (
	"fmt"
	"main/pkg/logger"
	"main/pkg/types"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

func GetConfig(path string) (*types.Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &types.Config{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	if err := defaults.Set(configStruct); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Error setting default config values")
	}

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
