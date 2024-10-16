package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (uc *UseCase) UserCreate(ctx context.Context, user models.UserCU) (int64, error) {
	return uc.srv.User.Create(ctx, user)
}
