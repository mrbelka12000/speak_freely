package usecase

import (
	"context"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/validate"
)

// UserCreate
func (uc *UseCase) UserCreate(ctx context.Context, user models.UserCU) (int64, map[string]validate.RequiredField, error) {

	id, err := uc.srv.User.Create(ctx, user)
	if err != nil {
		return 0, nil, fmt.Errorf("create user: %w", err)
	}

	return id, nil, nil
}

// UserUpdate
func (uc *UseCase) UserUpdate(ctx context.Context, pars models.UserGetPars, user models.UserCU) error {

	err := uc.srv.User.Update(ctx, pars, user)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// UserGet
func (uc *UseCase) UserGet(ctx context.Context, pars models.UserGetPars) (models.User, error) {
	user, err := uc.srv.User.Get(ctx, pars)
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
func (uc *UseCase) UserDelete(ctx context.Context, pars models.UserGetPars) error {
	return uc.srv.User.Delete(ctx, pars)
}
