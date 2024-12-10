package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/mrbelka12000/speak_freely/internal/usecase"
	"github.com/mrbelka12000/speak_freely/pkg/config"
)

type Cron struct {
	uc  *usecase.UseCase
	log *slog.Logger
	cfg config.Config
}

func NewCron(uc *usecase.UseCase, log *slog.Logger, cfg config.Config) *Cron {
	return &Cron{
		uc:  uc,
		log: log,
		cfg: cfg,
	}
}

// Start async worker
func (c *Cron) Start() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(c.cfg.GenerateInterval).Do(func() {
		c.addThemes()
	})

	s.StartAsync()
}

func (c *Cron) addThemes() {
	return
	for _, level := range []string{"A1", "A2", "B1", "B2", "C1"} {
		err := c.uc.ThemesGenerateWithAI(context.Background(), level)
		if err != nil {
			c.log.With("error", err).Error(fmt.Sprintf("failed to generate theme for %s", level))
		}
	}

	c.log.Info("all themes generated")
}
