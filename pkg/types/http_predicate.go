package types

import (
	"fmt"
	"net/http"
	"strconv"
)

type HTTPPredicate func(response *http.Response) error

func HTTPPredicateAlwaysPass() HTTPPredicate {
	return func(response *http.Response) error {
		return nil
	}
}

func HTTPPredicateCheckHeightAfter(prevHeight int64) HTTPPredicate {
	return func(response *http.Response) error {
		currentHeightHeader := response.Header.Get("Grpc-Metadata-X-Cosmos-Block-Height")

		// not returned height is ok
		if currentHeightHeader == "" {
			return nil
		}

		currentHeight, err := strconv.ParseInt(currentHeightHeader, 10, 64)
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
