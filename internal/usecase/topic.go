package usecase

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

func (uc *UseCase) TopicList(ctx context.Context, languageID int64) ([]models.Topic, error) {
	return uc.srv.Topic.List(ctx, languageID)
}
