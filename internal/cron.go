package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/usecase"
	"github.com/mrbelka12000/speak_freely/pkg/config"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type (
	Cron struct {
		uc       *usecase.UseCase
		notifier notifier
		log      *slog.Logger
		cfg      config.Config
	}

	notifier interface {
		NotifyTimeToPay(userID, chatID int64) error
	}
)

func NewCron(uc *usecase.UseCase, notifier notifier, log *slog.Logger, cfg config.Config) *Cron {
	return &Cron{
		uc:       uc,
		notifier: notifier,
		log:      log,
		cfg:      cfg,
	}
}

// Start async worker
func (c *Cron) Start() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(c.cfg.GenerateInterval).Do(func() {
		c.addThemes()
	})

	s.Every(1).Day().Do(func() {
		c.checkSubscriptionExpiration()
	})

	s.StartAsync()
}

func (c *Cron) addThemes() {
	for _, level := range []string{"A1", "A2", "B1", "B2", "C1"} {
		err := c.uc.ThemesGenerateWithAI(context.Background(), level)
		if err != nil {
			c.log.With("error", err).Error(fmt.Sprintf("failed to generate theme for %s", level))
		}
	}

	c.log.Info("all themes generated")
}

func (c *Cron) checkSubscriptionExpiration() {
	billingInfos, err := c.uc.BillingInfoList(context.Background())
	if err != nil {
		c.log.With("error", err).Error("failed to fetch billing infos")
		return
	}

	for _, bi := range billingInfos {
		if bi.DebitDate.Before(time.Now()) {
			err := c.notifier.NotifyTimeToPay(bi.UserID, bi.ChatID)
			if err != nil {
				c.log.With("error", err).Error("failed to notify time to pay")
				continue
			}

			err = c.uc.BillingInfoUpdate(context.Background(),
				bi.ID,
				models.BillingInfoCU{
					DebitDate: pointer.Of(time.Now().AddDate(0, 1, 1)),
				},
			)
			if err != nil {
				c.log.With("error", err).Error("failed to update billing info")
				continue
			}

			err = c.uc.UserUpdate(context.Background(),
				models.UserGetPars{
					ID: bi.UserID,
				},
				models.UserCU{
					RemainingTime: pointer.Of(int64(600)),
					Payed:         pointer.Of(false),
				},
			)
			if err != nil {
				c.log.With("error", err).Error("failed to update user subscription")
			}
		}
	}

	c.log.Info("all subscriptions checked")
}
