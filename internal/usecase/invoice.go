package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (uc *UseCase) ProceedInvoice(ctx context.Context, externalID, chatID int64) error {
	ctx, err := uc.tx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	err = uc.UserUpdate(ctx, models.UserGetPars{
		ExternalID: fmt.Sprint(externalID),
	}, models.UserCU{
		Payed: pointer.Of(true),
	})
	if err != nil {
		return fmt.Errorf("can not update user: %w", err)
	}

	user, err := uc.srv.User.Get(ctx, models.UserGetPars{
		ExternalID: fmt.Sprint(pointer.Value(&externalID)),
	})
	if err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	_, err = uc.BillingInfoCreate(ctx, models.BillingInfoCU{
		UserID:    pointer.Of(user.ID),
		ChatID:    pointer.Of(chatID),
		DebitDate: pointer.Of(time.Now().AddDate(0, 1, 0)),
	})
	if err != nil {
		return fmt.Errorf("can not create billing info: %w", err)
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
