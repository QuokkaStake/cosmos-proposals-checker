package types

import (
	"fmt"
	"main/pkg/utils"
	"net/http"
)

type HTTPPredicate func(response *http.Response) error

func HTTPPredicateAlwaysPass() HTTPPredicate {
	return func(response *http.Response) error {
		return nil
	}
}

func HTTPPredicateCheckHeightAfter(prevHeight int64) HTTPPredicate {
	return func(response *http.Response) error {
		currentHeight, err := utils.GetBlockHeightFromHeader(response.Header)
		if err != nil {
			return err
		}

		if prevHeight > currentHeight {
			return fmt.Errorf(
				"previous height (%d) is bigger than the current height (%d)",
				prevHeight,
				currentHeight,
			)
		}

		return nil
	}
}
