package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type language struct {
	db *sql.DB
}

func newLanguage(db *sql.DB) *language {
	return &language{
		db: db,
	}
}

func (l *language) Create(ctx context.Context, obj models.LanguageCU) error {
	_, err := l.db.ExecContext(ctx, "INSERT INTO languages (long_name, short_name) VALUES ($1,$2)", *obj.LongName, *obj.ShortName)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (l *language) Get(ctx context.Context, id int64) (obj models.Language, err error) {
	err = l.db.QueryRowContext(ctx, `SELECT id, short_name, long_name FROM languages WHERE id = $1`, id).
		Scan(&obj.ID, &obj.ShortName, &obj.LongName)
	if err != nil {
		return obj, fmt.Errorf("query: %w", err)
	}

	return obj, nil
}

func (l *language) GetByShortName(ctx context.Context, shortName string) (obj models.Language, err error) {
	err = l.db.QueryRowContext(ctx, `SELECT id, short_name, long_name FROM languages WHERE short_name = $1`, shortName).
		Scan(&obj.ID, &obj.ShortName, &obj.LongName)
	if err != nil {
		return obj, fmt.Errorf("query: %w", err)
	}

	return obj, nil
}

func (l *language) List(ctx context.Context) ([]models.Language, int, error) {
	rows, err := l.db.QueryContext(ctx, `SELECT id, short_name, long_name FROM languages`)
	if err != nil {
		return nil, 0, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var objs []models.Language

	for rows.Next() {
		var obj models.Language
		if err := rows.Scan(&obj.ID, &obj.ShortName, &obj.LongName); err != nil {
			return nil, 0, fmt.Errorf("scan: %w", err)
		}
		objs = append(objs, obj)
	}

	return objs, len(objs), nil
}
