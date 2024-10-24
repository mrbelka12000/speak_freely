package tgbot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

func (h *handler) getTopics(externalID int64) (empty tgbotapi.InlineKeyboardMarkup, err error) {
	ctx := context.Background()
	user, err := h.uc.UserGet(ctx, models.UserGet{ExternalID: fmt.Sprint(externalID)})
	if err != nil {
		return empty, fmt.Errorf("get user: %w", err)
	}

	themes, count, err := h.uc.ThemeList(ctx, models.ThemeListPars{
		LanguageID: pointer.Of(user.LanguageID),
		Random:     true,
		PaginationParams: models.PaginationParams{
			Limit: 10,
			Page:  1,
		},
	})
	if err != nil {
		return empty, fmt.Errorf("get themes: %w", err)
	}
	if count == 0 {
		return empty, nil
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, count)
	for _, theme := range themes {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: fmt.Sprintf("%s. %s. %s ", theme.Topic, theme.Level, theme.Question),
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseTheme,
					TC: &ThemeChoose{
						ID: theme.ID,
					},
				})),
			},
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	), nil
}
