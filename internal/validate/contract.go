package validate

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	userRepo interface {
		List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error)
	}
	langRepo interface {
		Get(ctx context.Context, id int64) (models.Language, error)
	}
)
