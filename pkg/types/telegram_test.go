package types

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTelegramResponseHasTextNotJson(t *testing.T) {
	t.Parallel()

	req := &http.Request{Body: io.NopCloser(strings.NewReader("not json"))}
	matcher := TelegramResponseHasText("text")
	require.False(t, matcher.Check(req))
}

func TestTelegramResponseHasTextDoesNotMatch(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	bytes, err := json.Marshal(TelegramResponse{Text: "text"})
	require.NoError(t, err)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(bytes)))}
	matcher := TelegramResponseHasText("wrong text")
	matcher.Check(req)
}
