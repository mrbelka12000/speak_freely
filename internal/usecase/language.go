package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (uc *UseCase) LanguageCreate(ctx context.Context, name string) error {
	return uc.srv.Language.Create(ctx, name)
}

func (uc *UseCase) LanguageGet(ctx context.Context, id int64) (models.Language, error) {
	return uc.srv.Language.Get(ctx, id)
}

func (uc *UseCase) LanguageList(ctx context.Context) ([]models.Language, int, error) {
	return uc.srv.Language.List(ctx)
}
