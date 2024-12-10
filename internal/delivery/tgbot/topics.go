package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func getTopics() tgbotapi.InlineKeyboardMarkup {
	preferences := []string{
		"Travelling", "Family", "Books", "Films", "Science", "Education", "Friends", "Social Media",
		"Work", "Cooking", "Personal Information", "Daily Routine", "Weather", "Food and Drink",
		"Hobbies", "Health", "Fitness", "Plans", "Cultural Differences",
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(preferences))
	for _, topic := range preferences {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text: fmt.Sprintf("%s", topic),
				CallbackData: pointer.Of(marshalCallbackData(CallbackData{
					Action: actionChooseTopic,
					TopC: &TopicChoose{
						Name: topic,
					},
				})),
			},
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
