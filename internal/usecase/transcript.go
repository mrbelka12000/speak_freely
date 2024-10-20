package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
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
