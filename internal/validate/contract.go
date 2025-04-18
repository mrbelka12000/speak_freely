package validate

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	userRepo interface {
		Get(ctx context.Context, pars models.UserGetPars) (models.User, error)
		List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error)
	}
	langRepo interface {
		Get(ctx context.Context, id int64) (models.Language, error)
	}
	fileRepo interface {
		Get(ctx context.Context, id int64) (models.File, error)
		GetByKey(ctx context.Context, key string) (models.File, error)
	}
	themeRepo interface {
		Get(ctx context.Context, id int64) (models.Theme, error)
	}
)
