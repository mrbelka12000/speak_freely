package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	PGURL          string `env:"PG_URL,required"`
	MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
	UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
	ServiceName    string `env:"SERVICE_NAME,required"`
	HTTPPort       string `env:"HTTP_PORT, default=8081"`
	AssemblyKey    string `env:"ASSEMBLY_KEY,required"`
	AIToken        string `env:"AI_TOKEN,required"`
}

func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	godotenv.Load()

	err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fill config: %w", err)
	}

	return cfg, nil
}
