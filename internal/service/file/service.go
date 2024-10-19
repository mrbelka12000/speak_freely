package file

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Service struct {
		r repo
	}
	repo interface {
		Create(ctx context.Context, obj models.FileCU) (int64, error)
		Get(ctx context.Context, id int64) (obj models.File, err error)
	}
)

func New(r repo) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, obj models.FileCU) (int64, error) {
	return s.r.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, id int64) (models.File, error) {
	return s.r.Get(ctx, id)
}
