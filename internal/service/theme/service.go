package theme

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, user models.ThemeCU) (int64, error)
		Get(ctx context.Context, id int64) (models.Theme, error)
		List(ctx context.Context, pars models.ThemeListPars) ([]models.Theme, int, error)
	}
)

func New(r repo) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Create(ctx context.Context, user models.ThemeCU) (int64, error) {
	return s.r.Create(ctx, user)
}

func (s *Service) Get(ctx context.Context, id int64) (models.Theme, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, pars models.ThemeListPars) ([]models.Theme, int, error) {
	return s.r.List(ctx, pars)
}
