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

func TestTelegramResponseHasTextWithMarkupNotJson(t *testing.T) {
	t.Parallel()

	req := &http.Request{Body: io.NopCloser(strings.NewReader("not json"))}
	matcher := TelegramResponseHasTextAndMarkup("text", TelegramInlineKeyboardResponse{})
	require.False(t, matcher.Check(req))
}

func TestTelegramResponseHasTextWithMarkupTextDoesNotMatch(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	bytes, err := json.Marshal(TelegramResponse{Text: "text"})
	require.NoError(t, err)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(bytes)))}
	matcher := TelegramResponseHasTextAndMarkup("wrong text", TelegramInlineKeyboardResponse{})
	matcher.Check(req)
}

func TestTelegramResponseHasTextWithMarkupTextKeyboardNotJSON(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	bytes, err := json.Marshal(TelegramResponse{Text: "text", ReplyMarkup: "not json"})
	require.NoError(t, err)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(bytes)))}
	matcher := TelegramResponseHasTextAndMarkup("text", TelegramInlineKeyboardResponse{})
	matcher.Check(req)
}

func TestTelegramResponseHasTextWithMarkupTextKeyboardDoesNotMatch(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	keyboardBytes, err := json.Marshal(TelegramInlineKeyboardResponse{
		InlineKeyboard: [][]TelegramInlineKeyboard{{
			{Unique: "unique"},
		}},
	})
	require.NoError(t, err)

	bytes, err := json.Marshal(TelegramResponse{Text: "text", ReplyMarkup: string(keyboardBytes)})
	require.NoError(t, err)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(bytes)))}
	matcher := TelegramResponseHasTextAndMarkup("text", TelegramInlineKeyboardResponse{
		InlineKeyboard: [][]TelegramInlineKeyboard{{}},
	})
	matcher.Check(req)
}
