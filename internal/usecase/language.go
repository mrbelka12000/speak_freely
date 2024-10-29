package usecase

import (
	"context"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/validate"
)

func (uc *UseCase) LanguageCreate(ctx context.Context, obj models.LanguageCU) (map[string]validate.RequiredField, error) {
	missed, err := uc.validator.ValidateLanguage(ctx, obj)
	if err != nil {
		return nil, fmt.Errorf("validate language error: %w", err)
	}
	if len(missed) != 0 {
		return missed, nil
	}
	err = uc.srv.Language.Create(ctx, obj)
	if err != nil {
		return nil, fmt.Errorf("create language error: %w", err)
	}

	return nil, nil
}

func (uc *UseCase) LanguageGet(ctx context.Context, id int64) (models.Language, error) {
	return uc.srv.Language.Get(ctx, id)
}

func (uc *UseCase) LanguageList(ctx context.Context) ([]models.Language, int, error) {
	return uc.srv.Language.List(ctx)
}
