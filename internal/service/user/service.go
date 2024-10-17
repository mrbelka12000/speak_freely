package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Service struct {
		r          repo
		bcryptCost int
	}

	repo interface {
		Create(ctx context.Context, user models.UserCU) (int64, error)
		Get(ctx context.Context, id int64) (models.User, error)
		Update(ctx context.Context, id int64, user models.UserCU) error
		List(ctx context.Context, pars models.UserPars) ([]models.User, int, error)
		Delete(ctx context.Context, id int64) error
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

// Create new user
func (s *Service) Create(ctx context.Context, user models.UserCU) (int64, error) {
	user.CreatedAt = time.Now().Unix()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), s.bcryptCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	*user.Password = string(hashedPassword)

	return s.r.Create(ctx, user)
}

// Get by id
func (s *Service) Get(ctx context.Context, id int64) (models.User, error) {
	return s.r.Get(ctx, id)
}

// Update user data
func (s *Service) Update(ctx context.Context, id int64, user models.UserCU) error {
	return s.r.Update(ctx, id, user)
}

// List
func (s *Service) List(ctx context.Context, pars models.UserPars) ([]models.User, int, error) {
	return s.r.List(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.r.Delete(ctx, id)
}

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
