package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
)

type (
	UseCase struct {
		srv       *service.Service
		tx        txBuilder
		validator *validate.Validator
	}

	txBuilder interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
)

func New(srv *service.Service, tx txBuilder, v *validate.Validator) *UseCase {
	return &UseCase{
		srv:       srv,
		tx:        tx,
		validator: v,
	}
}
