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

		cache cache
	}

	cache interface {
		Set(key string, value interface{}, dur time.Duration) error
		Get(key string) (string, bool)
		GetInt64(key string) (int64, bool)
		GetInt(key string) (int, bool)
	}
)

func Start(cfg config.Config, uc *usecase.UseCase, log *slog.Logger, cache cache) error {

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return fmt.Errorf("new bot: %w", err)
	}

	uCfg := tgbotapi.NewUpdate(0)
	uCfg.Timeout = 60

	h := handler{
		uc:    uc,
		log:   log,
		bot:   bot,
		ch:    bot.GetUpdatesChan(uCfg),
		cache: cache,
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

		//handle audio/video messages
		if id := getFileID(update.Message); id != "" {
			fileURL, err := h.getFileURL(id)
			if err != nil {
				h.log.With("error", err).Error("get file url")
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")))
				continue
			}

			id := update.Message.From.ID
			themeID, ok := h.cache.GetInt64(fmt.Sprintf("%d:theme", id))
			if !ok {
				h.log.With("error", err).Error("theme have not chosen")
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")))
				continue
			}

			msg, err := h.uc.TranscriptBuildFromURL(
				context.Background(),
				fileURL,
				themeID,
				id,
			)
			if err != nil {
				h.log.With("error", err).Error("save transcript")
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")))
				continue
			}

			h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)))
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
		case "themes":

			tMarkup, err := h.getThemes(msg.From.ID)
			if err != nil {
				h.log.With("error", err).Error("get topics")
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "something went wrong")))
				continue
			}

			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which theme do you want to discuss?")
			toSendMsg.ReplyMarkup = tMarkup
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
	case actionChooseTheme:
		tgUser := cb.From
		if tgUser == nil {
			h.log.Error("empty user in callback")
			return
		}

		theme, err := h.uc.ThemeGet(context.Background(), cbData.TC.ID)
		if err != nil {
			h.log.With("error", err).Error("get theme")
			return
		}

		h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, theme.Question)))
		err = h.cache.Set(fmt.Sprintf("%d:theme", tgUser.ID), cbData.TC.ID, 2*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set theme")
			return
		}
	}
}

func (h *handler) handleSendMessageError(_ tgbotapi.Message, err error) {
	if err != nil {
		h.log.With("error", err).Error("send message")
	}
}

func (h *handler) getFileURL(fileID string) (string, error) {
	url, err := h.bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", fmt.Errorf("get file direct url: %w", err)
	}

	return url, nil
}

func getFileID(msg *tgbotapi.Message) string {
	switch {
	case msg.Audio != nil:
		return msg.Audio.FileID
	case msg.Voice != nil:
		return msg.Voice.FileID
	case msg.Video != nil:
		return msg.Video.FileID
	case msg.VideoNote != nil:
		return msg.VideoNote.FileID
	default:
		return ""
	}
}
