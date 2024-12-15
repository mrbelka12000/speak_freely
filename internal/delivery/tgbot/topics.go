package tgbot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (h *Handler) getTopics(ctx context.Context, externalID int64) tgbotapi.InlineKeyboardMarkup {
	user, _ := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(externalID)})

	// Default language ID
	var langID int64 = 1
	if user.LanguageID != 0 {
		langID = user.LanguageID
	}

	topics, err := h.uc.TopicList(ctx, langID)
	if err != nil {
		h.log.With("error", err).Error("get topics")
		return tgbotapi.NewInlineKeyboardMarkup()
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(topics))
	for _, topic := range topics {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: fmt.Sprintf("%s", topic.Name),
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseTopic,
					TopC: &TopicChoose{
						ID: topic.ID,
					},
				})),
			},
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
