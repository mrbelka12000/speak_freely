package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/client/mail"
	handler "github.com/mrbelka12000/linguo_sphere_backend/internal/delivery/http/v1"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/repository"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/service"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/usecase"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/database"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/redis"
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

	rCache, err := redis.New(cfg)
	if err != nil {
		log.With("error", err).Error("failed to connect to redis")
		return
	}

	mailClient := mail.New(cfg)
	repo := repository.New(db)
	srv := service.New(repo)

	uc := usecase.New(
		srv,
		repo.Tx,
		validate.New(repo.User),
		mailClient,
		rCache,
		cfg.PublicURL,
		usecase.WithLogger(log),
	)

	h := handler.New(
		uc,
		handler.WithLogger(log),
	)

	r := mux.NewRouter()
	h.InitRoutes(r)

	s := server.New(r, cfg.HTTPPort)
	log.With("port", cfg.HTTPPort).Info("Starting HTTP server")
	s.Start() // non blocking

	gs := make(chan os.Signal, 1)
	signal.Notify(gs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-gs:
		log.Info(fmt.Sprintf("Received signal: %d", sig))
		log.Info("Server stopped properly")
		s.Stop()
		close(gs)
	case err := <-s.Ch():
		log.With("error", err).Error("Server stopped")
	}
}
