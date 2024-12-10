package tgbot

import "encoding/json"

type (
	CallbackData struct {
		Action action
		LangC  *LanguageChoose `json:"langc,omitempty"`
		ThemeC *ThemeChoose    `json:"themec,omitempty"`
		TopC   *TopicChoose    `json:"topc,omitempty"`
		LevelC *LevelChoose    `json:"levelc,omitempty"`
	}

	LanguageChoose struct {
		ID int64
	}

	ThemeChoose struct {
		ID int64 `json:"id"`
	}

	TopicChoose struct {
		Name string `json:"name"`
	}

	LevelChoose struct {
		Name string `json:"name"`
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
	actionChooseTheme    action = "choose_theme"
	actionChooseTopic    action = "choose_topic"
	actionChooseLevel    action = "choose_level"
)
