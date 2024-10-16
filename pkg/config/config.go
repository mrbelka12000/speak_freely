package config

type Config struct {
	PGURL          string `env:"PG_URL,required"`
	MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
	ServiceName    string `env:"SERVICE_NAME,required"`
	UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
	HTTPPort       string `env:"HTTP_PORT, default=5552"`
	AssemblyKey    string `env:"ASSEMBLY_KEY"`
	AIToken        string `env:"AI_TOKEN"`
}

func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	return Config{
		PGURL:          "",
		MigrationsPath: "migrations",
		ServiceName:    "linguo_sphere",
		UseMigrates:    false,
		HTTPPort:       "8080",
	}, nil
	//godotenv.Load()
	//
	//err = envconfig.Process(context.Background(), &cfg)
	//if err != nil {
	//	return cfg, fmt.Errorf("fill config: %w", err)
	//}

	return cfg, nil
}
