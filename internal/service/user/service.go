package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	Service struct {
		r          repo
		bcryptCost int
	}

	repo interface {
		Create(ctx context.Context, user models.UserCU) (int64, error)
		Get(ctx context.Context, obj models.UserGetPars) (models.User, error)
		Update(ctx context.Context, obj models.UserGetPars, user models.UserCU) error
		List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error)
		Delete(ctx context.Context, obj models.UserGetPars) error
	}
)

// New builder of user service
func New(r repo) *Service {
	return &Service{
		r:          r,
		bcryptCost: bcrypt.DefaultCost,
	}
}

// Create
func (s *Service) Create(ctx context.Context, user models.UserCU) (int64, error) {
	user.CreatedAt = time.Now().Unix()
	return s.r.Create(ctx, user)
}

// Get
func (s *Service) Get(ctx context.Context, pars models.UserGetPars) (models.User, error) {
	return s.r.Get(ctx, pars)
}

// Update
func (s *Service) Update(ctx context.Context, pars models.UserGetPars, user models.UserCU) error {
	return s.r.Update(ctx, pars, user)
}

// List
func (s *Service) List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error) {
	return s.r.List(ctx, pars)
}

// Delete
func (s *Service) Delete(ctx context.Context, pars models.UserGetPars) error {
	return s.r.Delete(ctx, pars)
}
