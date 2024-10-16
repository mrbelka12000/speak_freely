package validate

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type userRepo interface {
	List(ctx context.Context, pars models.UserPars) ([]models.User, int, error)
}
