package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type billingInfo struct {
	db *sql.DB
}

func newBillingInfo(db *sql.DB) *billingInfo {
	return &billingInfo{db: db}
}

func (b *billingInfo) Create(ctx context.Context, obj models.BillingInfoCU) (id int64, err error) {
	err = QueryRow(ctx, b.db, `
INSERT INTO billing_info(
	user_id, chat_id, debit_date              
) VALUES (
	$1, $2, $3
) RETURNING id`,
		pointer.Value(obj.UserID),
		pointer.Value(obj.ChatID),
		pointer.Value(obj.DebitDate)).Scan(&id)

	return id, err
}

func (b *billingInfo) Update(ctx context.Context, id int64, pars models.BillingInfoCU) error {

	queryUpdate := `
UPDATE billing_info
SET 
`
	queryWhere := "WHERE id = $1"
	args := []any{id}
	if pars.DebitDate != nil {
		args = append(args, *pars.DebitDate)
		queryUpdate += fmt.Sprintf(" debit_date = $%v ,", len(args))
	}

	queryUpdate = queryUpdate[:len(queryUpdate)-1]

	_, err := Exec(ctx, b.db, queryUpdate+queryWhere, args...)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (b *billingInfo) List(ctx context.Context) ([]models.BillingInfo, error) {
	rows, err := Query(ctx, b.db, `
SELECT id, user_id, debit_date, chat_id
FROM billing_info
`)
	if err != nil {
		return nil, fmt.Errorf("list rows: %w", err)
	}

	defer rows.Close()

	var result []models.BillingInfo

	for rows.Next() {
		var obj models.BillingInfo
		err := rows.Scan(
			&obj.ID,
			&obj.UserID,
			&obj.DebitDate,
			&obj.ChatID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, obj)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
