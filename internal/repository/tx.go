package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	Tx struct {
		db *sql.DB
	}

	txKey string
)

const (
	key txKey = "transaction_key"
)

func NewTx(db *sql.DB) *Tx {
	return &Tx{db: db}
}

func (t *Tx) Begin(ctx context.Context) (context.Context, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return context.WithValue(ctx, key, tx), nil
}

func (t *Tx) Commit(ctx context.Context) error {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get transaction from context: %w", err)
	}

	return tx.Commit()
}

func (t *Tx) Rollback(ctx context.Context) error {
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
