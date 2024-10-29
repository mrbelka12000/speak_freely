package language

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, obj models.LanguageCU) error
		Get(ctx context.Context, id int64) (obj models.Language, err error)
		GetByShortName(ctx context.Context, shortName string) (obj models.Language, err error)
		List(ctx context.Context) ([]models.Language, int, error)
	}
)

func New(r repo) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Create(ctx context.Context, obj models.LanguageCU) error {
	return s.r.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, id int64) (models.Language, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]models.Language, int, error) {
	return s.r.List(ctx)
}

func (s *Service) GetByShortName(ctx context.Context, shortName string) (obj models.Language, err error) {
	return s.r.GetByShortName(ctx, shortName)
}
