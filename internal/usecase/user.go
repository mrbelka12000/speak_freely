package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func (uc *UseCase) UserCreate(ctx context.Context, user models.UserCU) (int64, error) {
	return uc.srv.User.Create(ctx, user)
}

func (uc *UseCase) UserUpdate(ctx context.Context, id int64, user models.UserCU) error {
	return uc.srv.User.Update(ctx, id, user)
}

func (uc *UseCase) UserGet(ctx context.Context, id int64) (models.User, error) {
	return uc.srv.User.Get(ctx, id)
}

func (uc *UseCase) UserList(ctx context.Context, pars models.UserPars) ([]models.User, int, error) {
	return uc.srv.User.List(ctx, pars)
}

func (uc *UseCase) UserDelete(ctx context.Context, id int64) error {
	return uc.srv.User.Delete(ctx, id)
}

func (uc *UseCase) UserLogin(ctx context.Context, obj models.UserLogin) (int64, error) {
	return uc.srv.User.Login(ctx, obj)
}
