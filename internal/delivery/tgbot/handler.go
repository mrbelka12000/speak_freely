package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	lsb "github.com/mrbelka12000/linguo_sphere_backend"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

type (
	handler struct {
		uc  *usecase.UseCase
		log *slog.Logger

		bot *tgbotapi.BotAPI
		ch  tgbotapi.UpdatesChannel
	}

	cache interface {
		Set(key string, value interface{}, dur time.Duration) error
		GetInt64(key string) (int64, bool)
		GetInt(key string) (int, bool)
	}
)

func Start(cfg config.Config, uc *usecase.UseCase, log *slog.Logger) error {

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return fmt.Errorf("new bot: %w", err)
	}

	uCfg := tgbotapi.NewUpdate(0)
	uCfg.Timeout = 60

	h := handler{
		uc:  uc,
		log: log,
		bot: bot,
		ch:  bot.GetUpdatesChan(uCfg),
	}

	go h.handleUpdate()

	return nil
}

func (h *handler) handleUpdate() {
	for update := range h.ch {

		if update.CallbackQuery != nil {
			h.handleCallbacks(update.CallbackQuery)
			h.handleSendMessageError(h.bot.Send(
				tgbotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)),
			)
			continue
		}

		if update.Message == nil {
			continue
		}

		msg := update.Message
		switch msg.Command() {
		case "start":
			tgUser := msg.From

			_, err := h.uc.UserGet(context.Background(), models.UserGet{ExternalID: fmt.Sprint(tgUser.ID)})
			if err != nil {
				// user not exists, create
				obj := models.UserCU{
					FirstName:  pointer.Of(tgUser.FirstName),
					LastName:   pointer.Of(tgUser.LastName),
					Nickname:   pointer.Of(tgUser.UserName),
					AuthMethod: lsb.AuthMethodTG,
					LanguageID: pointer.Of(int64(1)),
					ExternalID: pointer.Of(fmt.Sprint(tgUser.ID)),
				}

				_, _, err = h.uc.UserCreate(context.Background(), obj)
				if err != nil {
					h.log.With("error", err).Error("failed to create user in tg")
					h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "something went wrong")))
					return
				}

				toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which language do you want to learn?")
				lMarkup, err := h.getLanguages()
				if err != nil {
					h.log.With("error", err).Error("get languages")
					h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "something went wrong")))
					continue
				}

				toSendMsg.ReplyMarkup = lMarkup
				h.handleSendMessageError(h.bot.Send(toSendMsg))
			} else {
				h.handleSendMessageError(
					h.bot.Send(
						tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Hello %s, let`s start to practice.", msg.From.UserName)),
					),
				)
			}

		case "language":
			lMarkup, err := h.getLanguages()
			if err != nil {
				h.log.With("error", err).Error("get languages")
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "something went wrong")))
				continue
			}

			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which language do you want to learn?")
			toSendMsg.ReplyMarkup = lMarkup
			h.handleSendMessageError(h.bot.Send(toSendMsg))
		}
	}
}

func (h *handler) handleCallbacks(cb *tgbotapi.CallbackQuery) {

	cbData, err := unmarshalCallbackData(cb.Data)
	if err != nil {
		h.log.With("error", err).Error("unmarshal callback data")
		return
	}

	switch cbData.Action {
	case actionChooseLanguage:
		tgUser := cb.From
		if tgUser == nil {
			h.log.Error("empty user in callback")
			return
		}

		_, err := h.uc.UserUpdate(context.Background(), models.UserGet{
			ExternalID: fmt.Sprint(tgUser.ID),
		}, models.UserCU{
			LanguageID: pointer.Of(cbData.LC.ID),
			AuthMethod: lsb.AuthMethodTG,
		})
		if err != nil {
			h.log.With("error", err).Error("update user")
			return
		}
	}
}

func (h *handler) handleSendMessageError(_ tgbotapi.Message, err error) {
	if err != nil {
		h.log.With("error", err).Error("send message")
	}
}
