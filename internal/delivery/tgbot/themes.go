package tgbot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (h *handler) getThemes(externalID int64) (empty tgbotapi.InlineKeyboardMarkup, err error) {
	ctx := context.Background()
	user, err := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(externalID)})
	if err != nil {
		return empty, fmt.Errorf("get user: %w", err)
	}

	var (
		level   *string
		topicID *int64
	)
	levelStr, ok := h.cache.Get(getLevelKey(externalID))
	if ok {
		level = &levelStr
	}

	topicStr, ok := h.cache.GetInt64(getTopicKey(externalID))
	if ok {
		topicID = &topicStr
	}

	themes, count, err := h.uc.ThemeList(ctx, models.ThemeListPars{
		LanguageID: pointer.Of(user.LanguageID),
		Level:      level,
		TopicID:    topicID,
		PaginationParams: models.PaginationParams{
			Limit: 10,
			Page:  1,
		},
	})
	if err != nil {
		return empty, fmt.Errorf("get themes: %w", err)
	}
	if count == 0 {
		return empty, fmt.Errorf("get themes: no themes found")
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, count)
	for _, theme := range themes {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: fmt.Sprintf("%s", theme.Question),
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseTheme,
					ThemeC: &ThemeChoose{
						ID: theme.ID,
					},
				})),
				SwitchInlineQueryCurrentChat: pointer.Of(fmt.Sprintf("%s", theme.Question)),
			},
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	), nil
}
