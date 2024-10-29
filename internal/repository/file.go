package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type (
	file struct {
		db *sql.DB
	}
)

func newFile(db *sql.DB) *file {
	return &file{
		db: db,
	}
}

func (f *file) Create(ctx context.Context, obj models.FileCU) (id int64, err error) {
	err = QueryRow(ctx, f.db, `
		INSERT INTO files(
			key
		)
		VALUES(
			$1
		) RETURNING id`,
		*obj.Key,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	return id, nil
}

func (f *file) Get(ctx context.Context, id int64) (file models.File, err error) {
	err = QueryRow(ctx, f.db, `
SELECT 
    id, 
    key
FROM files
WHERE id = $1`,
		id).Scan(
		&file.ID,
		&file.Key,
	)
	if err != nil {
		return file, fmt.Errorf("get file: %w", err)
	}

	return file, nil
}

func (f *file) GetByKey(ctx context.Context, key string) (file models.File, err error) {
	err = QueryRow(ctx, f.db, `
SELECT 
    id, 
    key
FROM files
WHERE key = $1`,
		key).Scan(
		&file.ID,
		&file.Key,
	)
	if err != nil {
		return file, fmt.Errorf("get file: %w", err)
	}

	return file, nil
}
