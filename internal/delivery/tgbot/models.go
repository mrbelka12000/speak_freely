package tgbot

import "encoding/json"

type (
	CallbackData struct {
		Action action          `json:"a"`
		LangC  *LanguageChoose `json:"lgc,omitempty"`
		ThemeC *ThemeChoose    `json:"thc,omitempty"`
		TopC   *TopicChoose    `json:"tpc,omitempty"`
		LevelC *LevelChoose    `json:"lvlc,omitempty"`
	}

	LanguageChoose struct {
		ID int64 `json:"id"`
	}

	ThemeChoose struct {
		ID int64 `json:"id"`
	}

	TopicChoose struct {
		ID int64 `json:"id"`
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
