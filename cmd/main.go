package main

import (
	"log/slog"
	"os"

	"github.com/gorilla/mux"

	handler "github.com/mrbelka12000/linguo_sphere_backend/internal/delivery/http/v1"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/database"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/server"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service_name", cfg.ServiceName)

	db, err := database.Connect(cfg)
	if err != nil {
		log.With("error", err).Error("failed to connect to database")
		return
	}
	defer db.Close()

	repo := repository.New(db)
	srv := service.New(repo)
	uc := usecase.New(srv, repo.Tx)
	h := handler.New(uc, validate.New(repo.User))
	r := mux.NewRouter()
	h.InitRoutes(r)

	s := server.New(r, cfg.HTTPPort)
	s.Start()

	log.With("port", cfg.HTTPPort).Info("Starting HTTP server")
	err = <-s.Ch()
	log.Error("Shutting down", err)
}
