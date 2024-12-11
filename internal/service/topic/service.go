package topic

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, obj models.Topic) error
		List(ctx context.Context, languageID int64) ([]models.Topic, error)
		Get(ctx context.Context, id int64) (obj models.Topic, err error)
	}
)

func New(r repo) *Service {
	return &Service{r: r}
}

func (s *Service) Get(ctx context.Context, id int64) (obj models.Topic, err error) {
	return s.r.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, obj models.Topic) error {
	return s.r.Create(ctx, obj)
}

func (s *Service) List(ctx context.Context, languageID int64) ([]models.Topic, error) {
	return s.r.List(ctx, languageID)
}
