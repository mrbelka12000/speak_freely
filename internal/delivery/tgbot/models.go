package tgbot

import "encoding/json"

type (
	CallbackData struct {
		Action action
		LC     *LanguageChoose `json:"lc,omitempty"`
	}

	LanguageChoose struct {
		ID int64
	}

	action string
)

func marshalCallbackData(cb CallbackData) string {
	body, _ := json.Marshal(cb)
	return string(body)
}

func unmarshalCallbackData(data string) (cb CallbackData, err error) {
	err = json.Unmarshal([]byte(data), &cb)
	return
}

const (
	actionChooseLanguage action = "choose_language"
)
