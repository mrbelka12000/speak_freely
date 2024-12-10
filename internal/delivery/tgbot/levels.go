package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func getLevels() tgbotapi.InlineKeyboardMarkup {
	levels := []string{"A1", "A2", "B1", "B2", "C1"}
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(levels))

	for _, level := range levels {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: fmt.Sprintf("%s", level),
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseLevel,
					LevelC: &LevelChoose{
						Name: level,
					},
				})),
			},
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
