package tgbot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (h *Handler) getThemes(externalID int64) (empty tgbotapi.InlineKeyboardMarkup, message string, err error) {
	ctx := context.Background()
	user, err := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(externalID)})
	if err != nil {
		return empty, "", fmt.Errorf("get user: %w", err)
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
		Random: true,
	})
	if err != nil {
		return empty, "", fmt.Errorf("get themes: %w", err)
	}
	if count == 0 {
		return empty, "", fmt.Errorf("get themes: no themes found")
	}

	var (
		response strings.Builder
		numbers  = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "1️⃣0️⃣"}
	)

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, count)
	button := make([]tgbotapi.InlineKeyboardButton, 0, 2)
	for i, theme := range themes {
		response.WriteString(fmt.Sprintf("%s: %s\n", numbers[i], theme.Question))
		button = append(button, tgbotapi.InlineKeyboardButton{
			Text: fmt.Sprintf("Choose %s", numbers[i]),
			CallbackData: pointer.Of(marshalCallbackData(CallbackData{
				Action: actionChooseTheme,
				ThemeC: &ThemeChoose{
					ID: theme.ID,
				},
			})),
			SwitchInlineQueryCurrentChat: pointer.Of(fmt.Sprintf("%s", theme.Question)),
		})
		if len(button) == 2 {
			buttons = append(buttons, button)
			button = make([]tgbotapi.InlineKeyboardButton, 0, 2)
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	), response.String(), nil
}
