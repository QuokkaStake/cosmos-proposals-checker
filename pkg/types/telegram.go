package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

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

func TelegramResponseHasBytesAndMarkup(text []byte, keyboard TelegramInlineKeyboardResponse) httpmock.Matcher {
	return TelegramResponseHasTextAndMarkup(string(text), keyboard)
}

func TelegramResponseHasTextAndMarkup(text string, keyboard TelegramInlineKeyboardResponse) httpmock.Matcher {
	return httpmock.NewMatcher("TelegramResponseHasTextAndMarkup",
		func(req *http.Request) bool {
			response := TelegramResponse{}

			err := json.NewDecoder(req.Body).Decode(&response)
			if err != nil {
				return false
			}

			if response.Text != text {
				panic(fmt.Sprintf("expected %q but got %q", response.Text, text))
			}

			var responseKeyboard TelegramInlineKeyboardResponse
			err = json.Unmarshal([]byte(response.ReplyMarkup), &responseKeyboard)
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(responseKeyboard, keyboard) {
				panic(fmt.Sprintf("expected keyboard %q but got %q", responseKeyboard, keyboard))
			}

			return true
		})
}
