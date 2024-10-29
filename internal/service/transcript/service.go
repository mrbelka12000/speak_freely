package transcript

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	Service struct {
		r repo
	}

	repo interface {
		Create(ctx context.Context, obj models.TranscriptCU) (id int64, err error)
		Get(ctx context.Context, id int64) (obj models.Transcript, err error)
		List(ctx context.Context, pars models.TranscriptListPars) ([]models.Transcript, int, error)
		Update(ctx context.Context, id int64, obj models.TranscriptCU) error
		Delete(ctx context.Context, id int64) error
	}
)

func New(r repo) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Create(ctx context.Context, obj models.TranscriptCU) (id int64, err error) {
	return s.r.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, id int64) (obj models.Transcript, err error) {
	return s.r.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, pars models.TranscriptListPars) ([]models.Transcript, int, error) {
	return s.r.List(ctx, pars)
}

func (s *Service) Update(ctx context.Context, id int64, obj models.TranscriptCU) error {
	return s.r.Update(ctx, id, obj)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.r.Delete(ctx, id)
}
