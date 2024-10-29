package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type (
	Service struct {
		r          repo
		bcryptCost int
	}

	repo interface {
		Create(ctx context.Context, user models.UserCU) (int64, error)
		Get(ctx context.Context, obj models.UserGet) (models.User, error)
		Update(ctx context.Context, obj models.UserGet, user models.UserCU) error
		List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error)
		Delete(ctx context.Context, obj models.UserGet) error
		FindByLogin(ctx context.Context, login string) (out models.User, err error)
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

	if user.AuthMethod == lsb.AuthMethodWeb {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pointer.Value(user.Password)), s.bcryptCost)
		if err != nil {
			return 0, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = pointer.Of(string(hashedPassword))
	}

	return s.r.Create(ctx, user)
}

// Get
func (s *Service) Get(ctx context.Context, pars models.UserGet) (models.User, error) {
	return s.r.Get(ctx, pars)
}

// Update
func (s *Service) Update(ctx context.Context, pars models.UserGet, user models.UserCU) error {
	return s.r.Update(ctx, pars, user)
}

// List
func (s *Service) List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error) {
	return s.r.List(ctx, pars)
}

// Delete
func (s *Service) Delete(ctx context.Context, pars models.UserGet) error {
	return s.r.Delete(ctx, pars)
}

// Login
func (s *Service) Login(ctx context.Context, obj models.UserLogin) (int64, error) {
	user, err := s.r.FindByLogin(ctx, obj.Login)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(obj.Password))
	if err != nil {
		return 0, fmt.Errorf("invalid password: %w", err)
	}

	return user.ID, nil
}
