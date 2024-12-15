package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

const (
	saveDuration = 30 * time.Minute
)

func (uc *UseCase) Conversation(
	ctx context.Context,
	fileURL string,
	externalUserID int64,
) (string, error) {
	ctx, err := uc.tx.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	user, err := uc.UserGet(ctx, models.UserGetPars{
		ExternalID: fmt.Sprint(externalUserID),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	var (
		prevBotAnswers, prevUserAnswer []string
	)

	if user.RemainingTime <= 0 && !user.Payed {
		return lsb.GetReachedLimitMessage(pointer.Value(user.Language).ShortName), nil
	}

	fileData, err := uc.transcriber.GetDataFromURL(ctx, fileURL, user.Language.ShortName)
	if err != nil {
		return "", fmt.Errorf("get user answer from message: %w", err)
	}

	botAnswerCache, ok := uc.getBotAnswer(externalUserID)
	if ok {
		prevBotAnswers = strings.Split(botAnswerCache, "---")
	}

	questionCache, ok := uc.getUserAnswer(externalUserID)
	if ok {
		prevUserAnswer = strings.Split(questionCache, "---")
	}

	dialog, err := uc.gen.Dialog(ctx, ai.DialogRequest{
		Text:      fileData.Text,
		Language:  user.Language.LongName,
		Questions: prevUserAnswer,
		Answers:   prevBotAnswers,
	})
	if err != nil {
		return "", fmt.Errorf("get dialog: %w", err)
	}

	err = uc.saveBotAnswer(externalUserID, dialog.Answer)
	if err != nil {
		uc.log.With("error", err).Error("can not save answer")
	}

	err = uc.saveUserAnswer(externalUserID, fileData.Text)
	if err != nil {
		uc.log.With("error", err).Error("can not save userAnswer")
	}

	err = uc.UserUpdate(ctx, models.UserGetPars{ID: user.ID}, models.UserCU{
		RemainingTime: pointer.Of(int64(fileData.AudioDuration) * -1),
	})
	if err != nil {
		return "", fmt.Errorf("update user info: %w", err)
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return dialog.Answer, nil
}

func (uc *UseCase) saveBotAnswer(externalUserID int64, text string) error {
	answer, ok := uc.cache.Get(GetBotAnswerKey(externalUserID))
	if ok {
		answer = fmt.Sprintf("%s---%s", answer, text)
	} else {
		answer = text
	}

	return uc.cache.Set(GetBotAnswerKey(externalUserID), answer, saveDuration)
}

func (uc *UseCase) getBotAnswer(externalUserID int64) (string, bool) {
	answer, ok := uc.cache.Get(GetBotAnswerKey(externalUserID))
	return answer, ok
}

func GetBotAnswerKey(externalUserID int64) string {
	return fmt.Sprintf("bot_answer_%d", externalUserID)
}

func (uc *UseCase) saveUserAnswer(externalUserID int64, text string) error {
	question, ok := uc.cache.Get(GetUserAnswerKey(externalUserID))
	if ok {
		question = fmt.Sprintf("%s---%s", question, text)
	} else {
		question = text
	}

	return uc.cache.Set(GetUserAnswerKey(externalUserID), question, saveDuration)
}

func (uc *UseCase) getUserAnswer(externalUserID int64) (string, bool) {
	question, ok := uc.cache.Get(GetUserAnswerKey(externalUserID))
	return question, ok
}

func GetUserAnswerKey(externalUserID int64) string {
	return fmt.Sprintf("user_answer_%d", externalUserID)
}
