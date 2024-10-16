package main

import (
	"log/slog"
	"net/http"
	"os"

	handler "github.com/mrbelka12000/linguo_sphere_backend/internal/delivery/http/v1"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/server"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service_name", cfg.ServiceName)

	repo := repository.New(nil)
	txBuilder := repository.NewTx(nil)
	srv := service.New(repo)
	uc := usecase.New(srv, txBuilder)

	h := handler.New(uc, validate.New(repo.User))
	mux := http.NewServeMux()
	h.InitRoutes(mux)

	s := server.New(mux, cfg.HTTPPort)
	s.Start()
	log.With("Port", cfg.HTTPPort).Info("Starting HTTP server")
	err = <-s.Ch()
	log.Error("Shutting down", err)
}
