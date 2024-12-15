package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/speak_freely/internal"
	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/client/assembly"
	handler "github.com/mrbelka12000/speak_freely/internal/delivery/http/v1"
	"github.com/mrbelka12000/speak_freely/internal/delivery/tgbot"
	"github.com/mrbelka12000/speak_freely/internal/repository"
	"github.com/mrbelka12000/speak_freely/internal/service"
	"github.com/mrbelka12000/speak_freely/internal/usecase"
	"github.com/mrbelka12000/speak_freely/internal/validate"
	"github.com/mrbelka12000/speak_freely/pkg/config"
	"github.com/mrbelka12000/speak_freely/pkg/database"
	"github.com/mrbelka12000/speak_freely/pkg/redis"
	"github.com/mrbelka12000/speak_freely/pkg/server"
	"github.com/mrbelka12000/speak_freely/pkg/storage/minio"
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

	minIOClient, err := minio.Connect(cfg)
	if err != nil {
		log.With("error", err).Error("failed to connect to minio")
		return
	}

	assemblyClient := assembly.New(cfg.AssemblyKey)
	aiClient := ai.NewClient(
		cfg.AIToken,
		ai.WithLogger(log),
	)
	repo := repository.New(db)
	srv := service.New(repo)

	uc := usecase.New(
		srv,
		repo.Tx,
		validate.New(repo.User, repo.Language, repo.File, repo.Theme),
		rCache,
		aiClient,
		minIOClient,
		assemblyClient,
		cfg.PublicURL,
		usecase.WithLogger(log),
	)

	tgHandler, err := tgbot.Start(cfg, uc, log, rCache) // non blocking
	if err != nil {
		log.With("error", err).Error("failed to start tgbot")
		return
	}

	internal.NewCron(uc, tgHandler, log, cfg).Start() // non blocking

	r := mux.NewRouter()

	handler.Init(
		uc,
		r,
		handler.WithLogger(log),
	)

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
