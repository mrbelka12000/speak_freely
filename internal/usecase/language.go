package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
)

func (uc *UseCase) LanguageCreate(ctx context.Context, obj models.LanguageCU) error {
	return uc.srv.Language.Create(ctx, obj)
}

func (uc *UseCase) LanguageGet(ctx context.Context, id int64) (models.Language, error) {
	return uc.srv.Language.Get(ctx, id)
}

func (uc *UseCase) LanguageList(ctx context.Context) ([]models.Language, int, error) {
	return uc.srv.Language.List(ctx)
}

func (uc *UseCase) LanguageValidate(ctx context.Context, obj models.LanguageCU) (map[string]validate.RequiredField, error) {
	return uc.validator.ValidateLanguage(ctx, obj)
}
