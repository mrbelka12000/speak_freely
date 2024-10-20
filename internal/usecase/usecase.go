package usecase

import (
	"html/template"
	"log/slog"
	"os"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
)

type (
	// UseCase
	UseCase struct {
		srv                  *service.Service
		tx                   txBuilder
		validator            *validate.Validator
		mailSender           mailSender
		cache                cache
		emailConfirmTemplate *template.Template
		gen                  generator
		storage              storage
		transcriber          transcriber

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
	ms mailSender,
	c cache,
	gen generator,
	s storage,
	t transcriber,
	publicURL string,
	opts ...opt,
) *UseCase {
	tmpl, err := template.ParseFiles("templates/email_confirmation.html")
	if err != nil {
		panic(err)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	uc := &UseCase{
		srv:                  srv,
		tx:                   tx,
		validator:            v,
		mailSender:           ms,
		cache:                c,
		emailConfirmTemplate: tmpl,
		gen:                  gen,
		storage:              s,
		transcriber:          t,

		log:       log,
		publicURL: publicURL,
	}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}
