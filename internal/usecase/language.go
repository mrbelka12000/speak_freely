package usecase

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

func (uc *UseCase) LanguageList(ctx context.Context) ([]models.Language, int, error) {
	return uc.srv.Language.List(ctx)
}

func (uc *UseCase) LanguageGet(ctx context.Context, id int64) (models.Language, error) {
	return uc.srv.Language.Get(ctx, id)
}
