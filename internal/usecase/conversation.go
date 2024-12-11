package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
)

const (
	saveDuration = 30 * time.Minute
)

func (uc *UseCase) Conversation(
	ctx context.Context,
	fileURL string,
	externalUserID int64,
) (string, error) {

	user, err := uc.UserGet(ctx, models.UserGet{
		ExternalID: fmt.Sprint(externalUserID),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	question, err := uc.transcriber.GetTextFromURL(ctx, fileURL, user.Language.ShortName)
	if err != nil {
		return "", fmt.Errorf("get question from message: %w", err)
	}

	var (
		prevAnswers, prevQuestions []string
	)
	answerCache, ok := uc.getAnswer(externalUserID)
	if ok {
		prevAnswers = strings.Split(answerCache, "---")
	}

	questionCache, ok := uc.getQuestion(externalUserID)
	if ok {
		prevQuestions = strings.Split(questionCache, "---")
	}

	dialog, err := uc.gen.Dialog(ctx, ai.DialogRequest{
		Text:      question,
		Language:  user.Language.LongName,
		Questions: prevQuestions,
		Answers:   prevAnswers,
	})
	if err != nil {
		return "", fmt.Errorf("get dialog: %w", err)
	}

	err = uc.saveAnswer(externalUserID, dialog.Answer)
	if err != nil {
		uc.log.With("error", err).Error("can not save answer")
	}

	err = uc.saveQuestion(externalUserID, question)
	if err != nil {
		uc.log.With("error", err).Error("can not save question")
	}

	return dialog.Answer, nil
}

func (uc *UseCase) saveAnswer(externalUserID int64, text string) error {
	answer, ok := uc.cache.Get(GetAnswerKey(externalUserID))
	if ok {
		answer = fmt.Sprintf("%s---%s", answer, text)
	} else {
		answer = text
	}

	return uc.cache.Set(GetAnswerKey(externalUserID), answer, saveDuration)
}

func (uc *UseCase) getAnswer(externalUserID int64) (string, bool) {
	answer, ok := uc.cache.Get(GetAnswerKey(externalUserID))
	return answer, ok
}

func GetAnswerKey(externalUserID int64) string {
	return fmt.Sprintf("answer_%d", externalUserID)
}

func (uc *UseCase) saveQuestion(externalUserID int64, text string) error {
	question, ok := uc.cache.Get(GetQuestionKey(externalUserID))
	if ok {
		question = fmt.Sprintf("%s---%s", question, text)
	} else {
		question = text
	}

	return uc.cache.Set(GetQuestionKey(externalUserID), question, saveDuration)
}

func (uc *UseCase) getQuestion(externalUserID int64) (string, bool) {
	question, ok := uc.cache.Get(GetQuestionKey(externalUserID))
	return question, ok
}

func GetQuestionKey(externalUserID int64) string {
	return fmt.Sprintf("question_%d", externalUserID)
}
