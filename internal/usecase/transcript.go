package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/validate"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (uc *UseCase) TranscriptBuild(
	ctx context.Context,
	file io.Reader,
	objectName,
	contentType string,
	fileSize int64,
	languageID int64,
	themeID int64,
	userID int64,
) (int64, map[string]validate.RequiredField, error) {
	ctx, err := uc.tx.Begin(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	fileID, err := uc.SaveFile(
		ctx,
		file,
		objectName,
		contentType,
		fileSize,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("save file: %w", err)
	}

	lang, err := uc.srv.Language.Get(ctx, languageID)
	if err != nil {
		return 0, nil, fmt.Errorf("get language: %w", err)
	}

	text, err := uc.transcriber.GetTextFromFile(ctx, file, lang.ShortName)
	if err != nil {
		return 0, nil, fmt.Errorf("get text: %w", err)
	}

	obj := models.TranscriptCU{
		Text:       pointer.Of(text),
		LanguageID: pointer.Of(languageID),
		UserID:     pointer.Of(userID),
		FileID:     pointer.Of(fileID),
		ThemeID:    pointer.Of(themeID),
	}

	missed, err := uc.validator.ValidateTranscript(ctx, obj, -1)
	if err != nil {
		return 0, nil, fmt.Errorf("transcript validation error: %w", err)
	}
	if len(missed) > 0 {
		return 0, missed, nil
	}

	id, err := uc.srv.Transcript.Create(ctx, obj)
	if err != nil {
		return 0, nil, fmt.Errorf("create transcript: %w", err)
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("commit transaction: %w", err)
	}

	return id, nil, nil
}

// TranscriptBuildFromURL telegram variant
func (uc *UseCase) TranscriptBuildFromURL(
	ctx context.Context,
	fileURL string,
	themeID int64,
	externalUserID int64,
) (msg string, err error) {
	ctx, err = uc.tx.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	user, err := uc.UserGet(ctx, models.UserGet{
		ExternalID: fmt.Sprint(externalUserID),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	text, err := uc.transcriber.GetTextFromURL(ctx, fileURL, user.Language.ShortName)
	if err != nil {
		return "", fmt.Errorf("get text from message: %w", err)
	}

	theme, err := uc.srv.Theme.Get(ctx, themeID)
	if err != nil {
		return "", fmt.Errorf("get theme: %w", err)
	}

	suggestions, err := uc.gen.GetSuggestions(ctx, ai.SuggestionRequest{
		Text:     text,
		Topic:    theme.Topic,
		Question: theme.Question,
		Language: user.Language.LongName,
	})
	if err != nil {
		return "", fmt.Errorf("get suggestions: %w", err)
	}
	suggestionRaw, _ := json.Marshal(suggestions)

	_, err = uc.srv.Transcript.Create(ctx, models.TranscriptCU{
		Text:       pointer.Of(text),
		LanguageID: pointer.Of(user.LanguageID),
		UserID:     pointer.Of(user.ID),
		ThemeID:    pointer.Of(themeID),
		Accuracy:   pointer.Of(suggestions.Accuracy),
		Suggestion: string(suggestionRaw),
	})
	if err != nil {
		return "", fmt.Errorf("create transcript: %w", err)
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return getSuggestionResponseTG(text, suggestions), nil
}

func (uc *UseCase) TranscriptGet(ctx context.Context, id int64, user models.User) (models.Transcript, error) {
	obj, err := uc.srv.Transcript.Get(ctx, id)
	if err != nil {
		return models.Transcript{}, err
	}

	if obj.UserID != user.ID {
		return models.Transcript{}, errors.New("user not permitted to get others people transcript")
	}

	return obj, nil
}

func (uc *UseCase) TranscriptList(ctx context.Context, pars models.TranscriptListPars) ([]models.Transcript, int, error) {
	return uc.srv.Transcript.List(ctx, pars)
}

func getSuggestionResponseTG(text string, s ai.SuggestionResponse) string {
	return fmt.Sprintf(`Accuracy: %v
	
Your text: %s

Corrected text: %s

Hints:
%v
`, s.Accuracy, text, s.Text, strings.Join(s.Hints, "\n"))
}
