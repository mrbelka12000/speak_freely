package language

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, name string) error
		Get(ctx context.Context, id int64) (obj models.Language, err error)
		List(ctx context.Context) ([]models.Language, int, error)
	}
)

func New(r repo) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Create(ctx context.Context, name string) error {
	return s.r.Create(ctx, name)
}

func (s *Service) Get(ctx context.Context, id int64) (models.Language, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]models.Language, int, error) {
	return s.r.List(ctx)
}
