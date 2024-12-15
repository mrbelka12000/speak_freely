package billing_info

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, obj models.BillingInfoCU) (int64, error)
		Update(ctx context.Context, id int64, obj models.BillingInfoCU) error
		List(ctx context.Context) ([]models.BillingInfo, error)
	}
)

func New(r repo) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, obj models.BillingInfoCU) (int64, error) {
	return s.r.Create(ctx, obj)
}

func (s *Service) Update(ctx context.Context, id int64, obj models.BillingInfoCU) error {
	return s.r.Update(ctx, id, obj)
}

func (s *Service) List(ctx context.Context) ([]models.BillingInfo, error) {
	return s.r.List(ctx)
}
