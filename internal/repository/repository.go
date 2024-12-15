package repository

import (
	"context"
	"database/sql"
)

type (
	// Repo
	Repo struct {
		User        *user
		Tx          *tx
		Language    *language
		Theme       *theme
		File        *file
		Transcript  *transcript
		Topic       *topic
		BillingInfo *billingInfo
	}

	contract interface {
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
)

// New
func New(db *sql.DB) *Repo {
	return &Repo{
		User:        newUser(db),
		Tx:          newTx(db),
		Language:    newLanguage(db),
		Theme:       newTheme(db),
		File:        newFile(db),
		Transcript:  newTranscript(db),
		Topic:       newTopic(db),
		BillingInfo: newBillingInfo(db),
	}
}

func getConnectionToDB(ctx context.Context, db *sql.DB) contract {
	tx, err := getTxFromContext(ctx)
	if err == nil {
		return tx
	}

	return db
}

// Exec custom exec function to implement transactions login
func Exec(ctx context.Context, db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	return getConnectionToDB(ctx, db).ExecContext(ctx, query, args...)
}

// Query custom query function to implement transactions login
func Query(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	return getConnectionToDB(ctx, db).QueryContext(ctx, query, args...)
}

// QueryRow custom queryRow function to implement transactions login
func QueryRow(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	return getConnectionToDB(ctx, db).QueryRowContext(ctx, query, args...)
}
