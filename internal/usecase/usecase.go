package usecase

import (
	"log/slog"
	"os"

	"github.com/mrbelka12000/speak_freely/internal/service"
	"github.com/mrbelka12000/speak_freely/internal/validate"
)

type (
	// UseCase
	UseCase struct {
		srv         *service.Service
		tx          txBuilder
		validator   *validate.Validator
		cache       cache
		gen         generator
		storage     storage
		transcriber transcriber

		log *slog.Logger

		publicURL string
		minIOURL  string
	}
)

// New
func New(
	srv *service.Service,
	tx txBuilder,
	v *validate.Validator,
	c cache,
	gen generator,
	s storage,
	t transcriber,
	publicURL string,
	opts ...opt,
) *UseCase {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	uc := &UseCase{
		srv:         srv,
		tx:          tx,
		validator:   v,
		cache:       c,
		gen:         gen,
		storage:     s,
		transcriber: t,

		log:       log,
		publicURL: publicURL,
	}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}
