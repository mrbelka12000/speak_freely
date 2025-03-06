package tgbot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/usecase"
	"github.com/mrbelka12000/speak_freely/pkg/config"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

const (
	smthWentWrong = "something went wrong"
)

type (
	Handler struct {
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

func Start(cfg config.Config, uc *usecase.UseCase, log *slog.Logger, cache cache) (*Handler, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("new bot: %w", err)
	}

	uCfg := tgbotapi.NewUpdate(0)
	uCfg.Timeout = 60

	h := Handler{
		uc:    uc,
		log:   log,
		bot:   bot,
		ch:    bot.GetUpdatesChan(uCfg),
		cache: cache,
	}

	go h.handleUpdate()

	return &h, nil
}

func (h *Handler) handleUpdate() {
	for update := range h.ch {

		ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)
		if update.Message != nil && update.Message.From != nil {
			_, err := h.uc.UserGet(ctx, models.UserGetPars{
				ExternalID: fmt.Sprint(update.Message.From.ID),
			})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					obj := models.UserCU{
						Nickname:   pointer.Of(update.Message.From.UserName),
						LanguageID: pointer.Of(int64(1)),
						ExternalID: pointer.Of(fmt.Sprint(update.Message.From.ID)),
					}

					_, _, err = h.uc.UserCreate(ctx, obj)
					if err != nil {
						h.log.With("error", err).Error("create user")
						continue
					}

					toSendMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Which language do you want to learn?")

					lMarkup, err := h.getLanguages()
					if err != nil {
						h.log.With("error", err).Error("get languages")
						continue
					}

					toSendMsg.ReplyMarkup = lMarkup
					h.handleSendMessageError(h.bot.Send(toSendMsg))
					continue
				}
			}
		}
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

		if update.PreCheckoutQuery != nil {
			err := h.handlePreCheckout(update.PreCheckoutQuery)
			if err != nil {
				h.log.With("error", err).Error("handle pre checkout")
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.SuccessfulPayment != nil {
			err := h.handleSuccessPayment(update.Message.SuccessfulPayment)
			if err != nil {
				h.log.With("error", err).Error("handle success payment")
				h.handleSendMessageError(
					h.bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						smthWentWrong),
					),
				)
				continue
			}
			h.handleSendMessageError(
				h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Thanks for payment<3",
				)),
			)
		}

		//handle audio/video messages
		if id := getFileID(update.Message); id != "" {
			fileURL, err := h.getFileURL(id)
			if err != nil {
				h.log.With("error", err).Error("get file url")
				h.handleSendMessageError(
					h.bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						smthWentWrong),
					),
				)
				continue
			}

			userID := update.Message.From.ID
			themeID, ok := h.cache.GetInt64(getThemeKey(userID))
			if !ok {
				answer, err := h.uc.Conversation(ctx, fileURL, userID)
				if err != nil {
					h.log.With("error", err).Error("get answer for conversation")
					h.handleSendMessageError(
						h.bot.Send(tgbotapi.NewMessage(
							update.Message.Chat.ID,
							smthWentWrong),
						),
					)
					continue
				}

				h.handleSendMessageError(
					h.bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						answer),
					),
				)
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
				h.handleSendMessageError(
					h.bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						smthWentWrong),
					),
				)
				continue
			}

			h.handleSendMessageError(
				h.bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					msg),
				),
			)
			continue
		}

		response, err := h.handleCommands(ctx, update.Message)
		if err != nil {
			h.log.With("error", err).Error("handle commands")
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						update.Message.Chat.ID,
						smthWentWrong,
					),
				),
			)
			continue
		}

		if response != "" {
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						update.Message.Chat.ID,
						response,
					),
				),
			)
		}

	}
}

func (h *Handler) handleCommands(ctx context.Context, msg *tgbotapi.Message) (string, error) {
	userID := msg.From.ID

	language, ok := h.cache.Get(getLanguageKey(userID))
	if !ok {
		user, err := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(userID)})
		if err != nil {
			h.log.With("error", err).Error("get user")
		}

		language = pointer.Value(user.Language).ShortName
		if language != "" {
			err = h.cache.Set(getLanguageKey(userID), language, 3*time.Hour)
			if err != nil {
				h.log.With("error", err).Error("set language to cache")
			}
		}
	}
	switch msg.Command() {

	case "start":
		_, err := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(userID)})
		if err != nil {
			// user not exists, create
			obj := models.UserCU{
				Nickname:   pointer.Of(msg.From.UserName),
				LanguageID: pointer.Of(int64(1)),
				ExternalID: pointer.Of(fmt.Sprint(userID)),
			}

			_, _, err = h.uc.UserCreate(ctx, obj)
			if err != nil {
				return smthWentWrong, fmt.Errorf("failed to create user: %w", err)
			}

			toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which language do you want to learn?")

			lMarkup, err := h.getLanguages()
			if err != nil {
				return smthWentWrong, fmt.Errorf("failed to get languages: %w", err)
			}
			toSendMsg.ReplyMarkup = lMarkup
			h.handleSendMessageError(h.bot.Send(toSendMsg))
		} else {
			return fmt.Sprintf(lsb.GetGreetingMessage(language), msg.From.UserName), nil
		}

	case "language":
		lMarkup, err := h.getLanguages()
		if err != nil {
			return smthWentWrong, fmt.Errorf("failed to get languages: %w", err)
		}

		toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Which language do you want to learn?")
		toSendMsg.ReplyMarkup = lMarkup
		h.handleSendMessageError(h.bot.Send(toSendMsg))

	case "themes":

		tMarkup, text, err := h.getThemes(msg.From.ID)
		if err != nil {
			h.log.With("error", err).Error("get themes")
			return lsb.GetNothingFindMessage(language), nil
		}

		toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, lsb.GetChooseThemeMessage(language)+"\n"+text)
		toSendMsg.ReplyMarkup = tMarkup
		h.handleSendMessageError(h.bot.Send(toSendMsg))

	case "reset":
		h.cache.Delete(usecase.GetUserAnswerKey(userID))
		h.cache.Delete(usecase.GetBotAnswerKey(userID))
		h.cache.Delete(getThemeKey(userID))
		h.cache.Delete(getTopicKey(userID))
		h.cache.Delete(getLevelKey(userID))

	case "topic":
		toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, lsb.GetChooseTopicMessage(language))
		toSendMsg.ReplyMarkup = h.getTopics(ctx, userID)
		h.handleSendMessageError(h.bot.Send(toSendMsg))

	case "level":
		toSendMsg := tgbotapi.NewMessage(msg.Chat.ID, lsb.GetChooseLevelMessage(language))
		toSendMsg.ReplyMarkup = getLevels()
		h.handleSendMessageError(h.bot.Send(toSendMsg))

	case "help":
		return lsb.GetFAQ(language), nil

	case "pay":
		user, _ := h.uc.UserGet(ctx, models.UserGetPars{ExternalID: fmt.Sprint(userID)})

		if user.Payed {
			return lsb.GetAlreadyPaidMessage(language), nil
		}

		err := h.sendInvoice(msg)
		if err != nil {
			return smthWentWrong, fmt.Errorf("failed to send invoice: %w", err)
		}

	case "apply":
		parts := strings.Split(msg.Text, " ")
		if len(parts) != 2 {
			return "invalid coupon provided", nil
		}

		return h.uc.ApplyCoupon(ctx, userID, parts[1])
	}
	return "", nil
}

func (h *Handler) handleCallbacks(ctx context.Context, cb *tgbotapi.CallbackQuery) {

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

		err := h.uc.UserUpdate(ctx, models.UserGetPars{
			ExternalID: fmt.Sprint(tgUser.ID),
		}, models.UserCU{
			LanguageID: pointer.Of(cbData.LangC.ID),
		})
		if err != nil {
			h.log.With("error", err).Error("update user")
			return
		}

		lang, err := h.uc.LanguageGet(ctx, cbData.LangC.ID)
		if err != nil {
			h.log.With("error", err).Error("get language")
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						cb.Message.Chat.ID,
						smthWentWrong,
					),
				),
			)
			return
		}

		err = h.cache.Set(getLanguageKey(tgUser.ID), lang.ShortName, 3*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set language to cache")
		}

	case actionChooseTheme:

		theme, err := h.uc.ThemeGet(ctx, cbData.ThemeC.ID)
		if err != nil {
			h.log.With("error", err).Error("get theme")
			return
		}

		err = h.cache.Set(getThemeKey(tgUser.ID), cbData.ThemeC.ID, 2*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set theme")
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						cb.Message.Chat.ID,
						smthWentWrong,
					),
				),
			)
			return
		}

		h.handleSendMessageError(
			h.bot.Send(
				tgbotapi.NewMessage(
					cb.Message.Chat.ID,
					theme.Question,
				),
			),
		)

	case actionChooseTopic:

		err = h.cache.Set(getTopicKey(tgUser.ID), cbData.TopC.ID, 2*time.Hour)
		if err != nil {
			h.log.With("error", err).Error("set topic")
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						cb.Message.Chat.ID,
						smthWentWrong,
					),
				),
			)

			return
		}

		topic, err := h.uc.TopicGet(ctx, cbData.TopC.ID)
		if err != nil {
			h.log.With("error", err).Error("get topic")
			return
		}

		h.handleSendMessageError(
			h.bot.Send(
				tgbotapi.NewMessage(
					cb.Message.Chat.ID,
					lsb.GetYourChooseMessage(h.getLanguageFromCache(cb.From.ID))+topic.Name),
			),
		)

	case actionChooseLevel:

		err = h.cache.Set(getLevelKey(tgUser.ID), cbData.LevelC.Name, 100*time.Hour)
		if err != nil {
			h.handleSendMessageError(
				h.bot.Send(
					tgbotapi.NewMessage(
						cb.Message.Chat.ID,
						smthWentWrong,
					),
				),
			)
			h.log.With("error", err).Error("set level")
			return
		}

		h.handleSendMessageError(
			h.bot.Send(
				tgbotapi.NewMessage(
					cb.Message.Chat.ID,
					lsb.GetYourChooseMessage(h.getLanguageFromCache(cb.From.ID))+cbData.LevelC.Name),
			),
		)
	}
}

func (h *Handler) handleSendMessageError(_ tgbotapi.Message, err error) {
	if err != nil {
		h.log.With("error", err).Error("send message")
	}
}

func (h *Handler) getFileURL(fileID string) (string, error) {
	url, err := h.bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", fmt.Errorf("get file direct url: %w", err)
	}

	return url, nil
}

func (h *Handler) getLanguageFromCache(userID int64) string {
	val, ok := h.cache.Get(getLanguageKey(userID))
	if ok {
		return val
	}

	h.log.Info("language cache miss")
	return ""
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

func getLanguageKey(id any) string {
	return fmt.Sprintf("%v:language", id)
}
