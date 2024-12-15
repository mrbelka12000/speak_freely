package tgbot

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type payload struct {
	UserID int64
	ChatID int64
}

func (h *Handler) handleSuccessPayment(payment *tgbotapi.SuccessfulPayment) error {
	var payload payload
	err := json.Unmarshal([]byte(payment.InvoicePayload), &payload)
	if err != nil {
		return err
	}

	err = h.uc.ProceedInvoice(context.Background(), payload.UserID, payload.ChatID)
	if err != nil {
		return fmt.Errorf("handle success payment: %w", err)
	}

	return nil
}

func (h *Handler) handlePreCheckout(data *tgbotapi.PreCheckoutQuery) error {
	pca := tgbotapi.PreCheckoutConfig{
		OK:                 true,
		PreCheckoutQueryID: data.ID,
	}

	{
		// Validate payload
		var payload payload
		err := json.Unmarshal([]byte(data.InvoicePayload), &payload)
		if err != nil {
			return err
		}
	}

	_, err := h.bot.Request(pca)
	if err != nil {
		return fmt.Errorf("answer: %w", err)
	}

	return nil
}

func (h *Handler) sendInvoice(msg *tgbotapi.Message) error {
	payload, err := makePayloadData(msg)
	if err != nil {
		return err
	}

	h.handleSendMessageError(h.bot.Send(tgbotapi.InvoiceConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: msg.Chat.ID,
		},
		Title:       "Unlimited access",
		Description: "Get unlimited access to the speaking, discussing, and improving language skills",
		Payload:     payload,
		Currency:    "XTR",
		Prices: []tgbotapi.LabeledPrice{
			{
				"50 stars",
				50,
			},
		},
		SuggestedTipAmounts: []int{},
	}))
	return nil
}

func makePayloadData(msg *tgbotapi.Message) (string, error) {
	payload := payload{
		UserID: msg.From.ID,
		ChatID: msg.Chat.ID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	return string(data), nil
}

func (h *Handler) NotifyTimeToPay(userID, chatID int64) error {
	if chatID == 0 {
		return fmt.Errorf("chatID is zero")
	}

	user, err := h.uc.UserGet(context.Background(), models.UserGetPars{
		ID: userID,
	})
	if err != nil {
		return fmt.Errorf("can not get user: %w", err)
	}

	h.handleSendMessageError(
		h.bot.Send(
			tgbotapi.NewMessage(
				chatID,
				lsb.GetTimeToPayMessage(pointer.Value(user.Language).ShortName),
			),
		),
	)

	return nil
}
