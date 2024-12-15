package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

const (
	envCouponAlem  = "COUPON_ALEM"
	envCouponOther = "COUPON_OTHER"
)

func (uc *UseCase) ApplyCoupon(ctx context.Context, externalID int64, coupon string) (string, error) {
	var valid bool

	secret, ok := os.LookupEnv(envCouponAlem)
	if ok && secret == coupon {
		valid = true
	}

	secret, ok = os.LookupEnv(envCouponOther)
	if ok && secret == coupon {
		valid = true
	}

	if !valid {
		return "invalid coupon provided", nil
	}

	ctx, err := uc.tx.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	user, err := uc.srv.User.Get(ctx, models.UserGetPars{
		ExternalID: fmt.Sprint(externalID),
	})
	if err != nil {
		return "", err
	}

	if user.IsRedeemUsed {
		return "coupon already used", nil
	}

	err = uc.UserUpdate(ctx,
		models.UserGetPars{
			ExternalID: fmt.Sprint(externalID),
		},
		models.UserCU{
			RemainingTime: pointer.Of(int64(3 * 60)),
			IsRedeemUsed:  pointer.Of(true),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return fmt.Sprintf("coupon successfully used, added %d hours", 3), nil
}
