package usecase

import (
	"context"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
)

type (
	UseCase struct {
		srv *service.Service
		tx  txBuilder
	}

	txBuilder interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
)

func New(srv *service.Service, tx txBuilder) *UseCase {
	return &UseCase{
		srv: srv,
		tx:  tx,
	}
}
