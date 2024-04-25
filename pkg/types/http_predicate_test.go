package types

import (
	"main/pkg/constants"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPPredicateAlwaysPass(t *testing.T) {
	t.Parallel()

	predicate := HTTPPredicateAlwaysPass()
	require.NoError(t, predicate(&http.Response{}))
}

func TestHTTPPredicateCheckHeightAfterErrorParsing(t *testing.T) {
	t.Parallel()

	predicate := HTTPPredicateCheckHeightAfter(1)

	header := http.Header{
		constants.HeaderBlockHeight: []string{"invalid"},
	}
	request := &http.Response{Header: header}
	require.Error(t, predicate(request))
}

func TestHTTPPredicateCheckHeightAfterOlderBlock(t *testing.T) {
	t.Parallel()

	predicate := HTTPPredicateCheckHeightAfter(100)

	header := http.Header{
		constants.HeaderBlockHeight: []string{"1"},
	}
	request := &http.Response{Header: header}
	require.Error(t, predicate(request))
}

func TestHTTPPredicateCheckHeightPass(t *testing.T) {
	t.Parallel()

	predicate := HTTPPredicateCheckHeightAfter(100)

	header := http.Header{
		constants.HeaderBlockHeight: []string{"200"},
	}
	request := &http.Response{Header: header}
	require.NoError(t, predicate(request))
}
