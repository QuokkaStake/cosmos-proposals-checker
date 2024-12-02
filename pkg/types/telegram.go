package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

type TelegramResponse struct {
	ChatID      string `json:"chat_id"`
	Text        string `json:"text"`
	ReplyMarkup string `json:"reply_markup"`
}

type TelegramInlineKeyboardResponse struct {
	InlineKeyboard [][]TelegramInlineKeyboard `json:"inline_keyboard"`
}

type TelegramInlineKeyboard struct {
	Unique       string `json:"unique"`
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func TelegramResponseHasBytes(text []byte) httpmock.Matcher {
	return TelegramResponseHasText(string(text))
}

func TelegramResponseHasText(text string) httpmock.Matcher {
	return httpmock.NewMatcher("TelegramResponseHasText",
		func(req *http.Request) bool {
			response := TelegramResponse{}
			err := json.NewDecoder(req.Body).Decode(&response)
			if err != nil {
				return false
			}

			if response.Text != text {
				panic(fmt.Sprintf("expected %q but got %q", response.Text, text))
			}

			return true
		})
}
