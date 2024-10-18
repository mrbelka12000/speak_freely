package usecase

import (
	"context"
	"fmt"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
)

// UserCreate
func (uc *UseCase) UserCreate(ctx context.Context, user models.UserCU) (int64, error) {
	id, err := uc.srv.User.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	go uc.sendConfirmationEmail(context.WithoutCancel(ctx), id)
	return id, nil
}

// UserUpdate
func (uc *UseCase) UserUpdate(ctx context.Context, id int64, user models.UserCU) error {
	return uc.srv.User.Update(ctx, id, user)
}

// UserGet
func (uc *UseCase) UserGet(ctx context.Context, id int64) (models.User, error) {
	user, err := uc.srv.User.Get(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}

	lang, err := uc.srv.Language.Get(ctx, user.LanguageID)
	if err != nil {
		return models.User{}, fmt.Errorf("get language: %w", err)
	}
	user.Language = &lang

	return user, nil
}

// UserList
func (uc *UseCase) UserList(ctx context.Context, pars models.UserListPars) ([]models.User, int, error) {
	return uc.srv.User.List(ctx, pars)
}

// UserDelete
func (uc *UseCase) UserDelete(ctx context.Context, id int64) error {
	return uc.srv.User.Delete(ctx, id)
}

// UserLogin
func (uc *UseCase) UserLogin(ctx context.Context, obj models.UserLogin) (int64, error) {
	return uc.srv.User.Login(ctx, obj)
}

// UserCUValidate
func (uc *UseCase) UserCUValidate(ctx context.Context, user models.UserCU, id int64) (map[string]validate.RequiredField, error) {
	return uc.validator.ValidateUser(ctx, user, id)
}
