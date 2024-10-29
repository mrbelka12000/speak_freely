package usecase

import (
	"context"
	"fmt"

	lsb "github.com/mrbelka12000/speak_freely"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/validate"
)

// UserCreate
func (uc *UseCase) UserCreate(ctx context.Context, user models.UserCU) (int64, map[string]validate.RequiredField, error) {
	missed, err := uc.validator.ValidateUser(ctx, user, -1)
	if err != nil {
		return 0, nil, fmt.Errorf("validate user: %w", err)
	}
	if len(missed) > 0 {
		return 0, missed, nil
	}

	id, err := uc.srv.User.Create(ctx, user)
	if err != nil {
		return 0, nil, fmt.Errorf("create user: %w", err)
	}

	if user.AuthMethod == lsb.AuthMethodWeb {
		go uc.sendConfirmationEmail(context.WithoutCancel(ctx), id)
	}
	return id, nil, nil
}

// UserUpdate
func (uc *UseCase) UserUpdate(ctx context.Context, pars models.UserGet, user models.UserCU) (map[string]validate.RequiredField, error) {
	missed, err := uc.validator.ValidateUser(ctx, user, pars.ID)
	if err != nil {
		return nil, fmt.Errorf("validate user: %w", err)
	}
	if len(missed) > 0 {
		return missed, nil
	}

	err = uc.srv.User.Update(ctx, pars, user)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return nil, nil
}

// UserGet
func (uc *UseCase) UserGet(ctx context.Context, pars models.UserGet) (models.User, error) {
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
func (uc *UseCase) UserDelete(ctx context.Context, pars models.UserGet) error {
	return uc.srv.User.Delete(ctx, pars)
}

// UserLogin
func (uc *UseCase) UserLogin(ctx context.Context, obj models.UserLogin) (int64, error) {
	return uc.srv.User.Login(ctx, obj)
}
