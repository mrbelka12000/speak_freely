package usecase

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

func (uc *UseCase) BillingInfoCreate(ctx context.Context, obj models.BillingInfoCU) (int64, error) {
	return uc.srv.BillingInfo.Create(ctx, obj)
}

func (uc *UseCase) BillingInfoUpdate(ctx context.Context, id int64, obj models.BillingInfoCU) error {
	return uc.srv.BillingInfo.Update(ctx, id, obj)
}

func (uc *UseCase) BillingInfoList(ctx context.Context) ([]models.BillingInfo, error) {
	return uc.srv.BillingInfo.List(ctx)
}
