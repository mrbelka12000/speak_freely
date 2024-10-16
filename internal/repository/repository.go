package repository

import (
	"context"
	"database/sql"
)

type (
	Repo struct {
		User *user
	}

	contract interface {
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
)

func New(db *sql.DB) *Repo {
	return &Repo{
		User: newUser(db),
	}
}

func getConnectionToDB(ctx context.Context, db *sql.DB) contract {
	tx, err := getTxFromContext(ctx)
	if err == nil {
		return tx
	}

	return db
}

func Exec(ctx context.Context, db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	return getConnectionToDB(ctx, db).ExecContext(ctx, query, args...)
}

func Query(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	return getConnectionToDB(ctx, db).QueryContext(ctx, query, args...)
}

func QueryRow(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	return getConnectionToDB(ctx, db).QueryRowContext(ctx, query, args...)
}
