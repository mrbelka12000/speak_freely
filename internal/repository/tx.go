package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	tx struct {
		db *sql.DB
	}

	txKey string
)

const (
	key txKey = "transaction_key"
)

func newTx(db *sql.DB) *tx {
	return &tx{db: db}
}

func (t *tx) Begin(ctx context.Context) (context.Context, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return context.WithValue(ctx, key, tx), nil
}

func (t *tx) Commit(ctx context.Context) error {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get transaction from context: %w", err)
	}

	return tx.Commit()
}

func (t *tx) Rollback(ctx context.Context) error {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get transaction from context: %w", err)
	}

	return tx.Rollback()
}

func getTxFromContext(ctx context.Context) (*sql.Tx, error) {
	txAny := ctx.Value(key)
	if txAny == nil {
		return nil, fmt.Errorf("transaction context not exists")
	}

	tx, ok := txAny.(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("tx from context is not of type *sql.Tx")
	}

	return tx, nil
}
