package user

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, user models.UserCU) (int64, error)
		Get(ctx context.Context, id int64) (models.User, error)
		Update(ctx context.Context, id int64, user models.UserCU) error
		List(ctx context.Context, pars models.UserPars) ([]models.User, int, error)
		Delete(ctx context.Context, id int64) error
	}
)

func New(r repo) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, user models.UserCU) (int64, error) {
	return s.r.Create(ctx, user)
}

func (s *Service) Get(ctx context.Context, id int64) (models.User, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, user models.UserCU) error {
	return s.r.Update(ctx, id, user)
}

func (s *Service) List(ctx context.Context, pars models.UserPars) ([]models.User, int, error) {
	return s.r.List(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.r.Delete(ctx, id)
}
