package usecase

import (
	"context"
	"time"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/client/mail"
)

type (
	mailSender interface {
		Send(req mail.Request) error
	}

	txBuilder interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}

	cache interface {
		Set(key string, value interface{}, dur time.Duration) error
		Get(key string) (string, bool)
		GetInt64(key string) (int64, bool)
		Delete(key string)
	}
)