package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/mrbelka12000/speak_freely/pkg/config"
)

// Connect
func Connect(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PGURL)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	if cfg.UseMigrates {
		err = useMigrates(db, cfg)
		if err != nil {
			return nil, fmt.Errorf("migrates: %w", err)
		}
	}

	return db, nil
}

func useMigrates(db *sql.DB, cfg config.Config) error {
	var id int
	err := db.QueryRow(`
		SELECT id FROM users LIMIT 1
`).Scan(&id)
	if err == nil {
		return nil
	}

	dir, err := os.ReadDir(cfg.MigrationsPath)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if strings.Contains(file.Name(), "down") {
			continue
		}

		query, _ := os.ReadFile(cfg.MigrationsPath + file.Name())

		_, err = db.Exec(string(query))
		if err != nil {
			return err
		}
	}
	return nil
}
