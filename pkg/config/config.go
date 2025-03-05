package config

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	// Config of service
	Config struct {
		InstanceConfig
		DBConfig
		ClientsConfig
		SMTPConfig
		RedisConfig
		MinIOConfig
		TelegramConfig
		CronConfig
	}

	InstanceConfig struct {
		ServiceName string `env:"SERVICE_NAME,required"`
		HTTPPort    string `env:"HTTP_PORT, default=8081"`
		PublicURL   string `env:"PUBLIC_URL,required"`
	}

	DBConfig struct {
		PGURL          string `env:"PG_URL,required"`
		MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
		UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
	}

	ClientsConfig struct {
		AIToken     string `env:"AI_TOKEN,required"`
		AssemblyKey string `env:"ASSEMBLY_KEY,required"`
	}

	SMTPConfig struct {
		SMTPHost     string `env:"SMTP_HOST, default=smtp.gmail.com"`
		SMTPPort     int    `env:"SMTP_PORT, default=587"`
		SenderEmail  string `env:"SENDER_EMAIL,required"`
		SMTPPassword string `env:"SENDER_PASSWORD,required"`
	}

	RedisConfig struct {
		RedisAddr string `env:"REDIS_ADDR,required"`
	}

	MinIOConfig struct {
		MinIOAddr      string `env:"MINIO_ADDR,required"`
		MinIOBucket    string `env:"MINIO_BUCKET,default=linguo_sphere"`
		MinIOAccessKey string `env:"MINIO_ACCESS_KEY,required"`
		MinIOSecretKey string `env:"MINIO_SECRET_KEY,required"`
	}

	TelegramConfig struct {
		BotToken string `env:"BOT_TOKEN,required"`
	}

	CronConfig struct {
		GenerateInterval time.Duration `env:"GENERATE_INTERVAL,default=12h"`
	}
)

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
