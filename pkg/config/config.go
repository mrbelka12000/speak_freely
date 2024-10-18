package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config of service
type Config struct {
	PGURL          string `env:"PG_URL,required"`
	MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
	UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
	ServiceName    string `env:"SERVICE_NAME,required"`
	HTTPPort       string `env:"HTTP_PORT, default=8081"`
	AssemblyKey    string `env:"ASSEMBLY_KEY,required"`
	AIToken        string `env:"AI_TOKEN,required"`
	SMTPHost       string `env:"SMTP_HOST, default=smtp.gmail.com"`
	SMTPPort       int    `env:"SMTP_PORT, default=587"`
	SenderEmail    string `env:"SENDER_EMAIL,required"`
	SMTPPassword   string `env:"SENDER_PASSWORD,required"`
	PublicURL      string `env:"PUBLIC_URL,required"`
	RedisAddr      string `env:"REDIS_ADDR,required"`
	RedisPassword  string `env:"REDIS_PASSWORD,required"`
}

// Get
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
