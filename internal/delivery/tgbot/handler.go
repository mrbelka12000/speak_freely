package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/usecase"
	"github.com/mrbelka12000/speak_freely/pkg/config"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
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
		Delete(key string)
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
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

		if update.CallbackQuery != nil {
			h.handleCallbacks(ctx, update.CallbackQuery)
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

			userID := update.Message.From.ID
			themeID, ok := h.cache.GetInt64(fmt.Sprintf("%d:theme", userID))
			if !ok {
				answer, err := h.uc.Conversation(ctx, fileURL, userID)
				if err != nil {
					h.log.With("error", err).Error("get answer for conversation")
					h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")))
					continue
				}

				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, answer)))
				continue
			}

			msg, err := h.uc.TranscriptBuildFromURL(
				ctx,
				fileURL,
				themeID,
				userID,
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
		tgUser := msg.From
		switch msg.Command() {
		case "start":
			_, err := h.uc.UserGet(ctx, models.UserGet{ExternalID: fmt.Sprint(tgUser.ID)})
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

				_, _, err = h.uc.UserCreate(ctx, obj)
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
				h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "По вашим критериям ничего не найдено, попробуйте выбрать другую тему.")))
				continue
			}

			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which theme do you want to discuss?")
			toSendMsg.ReplyMarkup = tMarkup
			h.handleSendMessageError(h.bot.Send(toSendMsg))
		case "conversation":
			h.cache.Delete(getThemeKey(update.Message.From.ID))
			h.cache.Delete(getTopicKey(update.Message.From.ID))
			h.cache.Delete(getLevelKey(update.Message.From.ID))
			h.cache.Delete(usecase.GetQuestionKey(update.Message.From.ID))
			h.cache.Delete(usecase.GetAnswerKey(update.Message.From.ID))
		case "topic":
			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which topic do you want to discuss?")
			toSendMsg.ReplyMarkup = h.getTopics(ctx, tgUser.ID)
			h.handleSendMessageError(h.bot.Send(toSendMsg))

		case "level":
			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which level are you ?")
			toSendMsg.ReplyMarkup = getLevels()
			h.handleSendMessageError(h.bot.Send(toSendMsg))
		case "help":
			user, _ := h.uc.UserGet(ctx, models.UserGet{ExternalID: fmt.Sprint(tgUser.ID)})
			faq := lsb.GetFAQ(pointer.Value(user.Language).ShortName)
			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, faq)
			h.handleSendMessageError(h.bot.Send(toSendMsg))
		}
	}
}

func (h *handler) handleCallbacks(ctx context.Context, cb *tgbotapi.CallbackQuery) {

	cbData, err := unmarshalCallbackData(cb.Data)
	if err != nil {
		h.log.With("error", err).Error("unmarshal callback data")
		return
	}

	tgUser := cb.From
	if tgUser == nil {
		h.log.Error("empty user in callback")
		return
	}
	switch cbData.Action {
	case actionChooseLanguage:

		_, err := h.uc.UserUpdate(ctx, models.UserGet{
			ExternalID: fmt.Sprint(tgUser.ID),
		}, models.UserCU{
			LanguageID: pointer.Of(cbData.LangC.ID),
			AuthMethod: lsb.AuthMethodTG,
		})
		if err != nil {
			h.log.With("error", err).Error("update user")
			return
		}
	case actionChooseTheme:

		theme, err := h.uc.ThemeGet(ctx, cbData.ThemeC.ID)
		if err != nil {
			h.log.With("error", err).Error("get theme")
			return
		}

		h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, theme.Question)))
		err = h.cache.Set(getThemeKey(tgUser.ID), cbData.ThemeC.ID, 2*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set theme")
			return
		}
	case actionChooseTopic:

		err = h.cache.Set(getTopicKey(tgUser.ID), cbData.TopC.ID, 2*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set topic")
			return
		}

	case actionChooseLevel:
		h.handleSendMessageError(h.bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Your choice: "+cbData.LevelC.Name)))
		err = h.cache.Set(getLevelKey(tgUser.ID), cbData.LevelC.Name, 100*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set level")
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
	case msg.VideoNote != nil:
		return msg.VideoNote.FileID
	default:
		return ""
	}
}

func getThemeKey(id int64) string {
	return fmt.Sprintf("%d:theme", id)
}

func getTopicKey(id any) string {
	return fmt.Sprintf("%v:topic", id)
}

func getLevelKey(id any) string {
	return fmt.Sprintf("%v:level", id)
}
