package v1

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mrbelka12000/speak_freely/internal/usecase"
)

type (
	Handler struct {
		uc      *usecase.UseCase
		log     *slog.Logger
		decoder *schema.Decoder
	}
)

// Init
func Init(uc *usecase.UseCase, router *mux.Router, opts ...opt) {

	h := &Handler{
		uc:  uc,
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	for _, opt := range opts {
		opt(h)
	}

	h.initRoutes(router)
	return
}

// InitRoutes
func (h *Handler) initRoutes(r *mux.Router) {
	r.Use(h.recovery)
	r.Use(h.cors)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Handle("/metrics", promhttp.Handler())
}
