package tgbot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (h *Handler) getLanguages() (empty tgbotapi.InlineKeyboardMarkup, err error) {
	languages, count, err := h.uc.LanguageList(context.Background())
	if err != nil {
		return empty, fmt.Errorf("can not get languages: %w", err)
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, count)
	for _, l := range languages {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: l.LongName,
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseLanguage,
					LangC:  &LanguageChoose{ID: l.ID},
				})),
			},
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	), nil
}
