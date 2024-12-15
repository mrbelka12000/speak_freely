package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type user struct {
	db *sql.DB
}

func newUser(db *sql.DB) *user {
	return &user{
		db: db,
	}
}

// Create
func (u *user) Create(ctx context.Context, user models.UserCU) (id int64, err error) {
	err = QueryRow(ctx, u.db, `
		INSERT INTO users(
			nickname,
			language_id,
			created_at,
		    external_id
			) 
		VALUES(
			$1,
			$2,
			$3,
			$4
		) RETURNING id`,
		pointer.Value(user.Nickname),
		pointer.Value(user.LanguageID),
		user.CreatedAt,
		pointer.Value(user.ExternalID),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

// Get
func (u *user) Get(ctx context.Context, obj models.UserGetPars) (user models.User, err error) {
	var (
		queryWhere string
		arg        any
	)
	if obj.ID != 0 {
		queryWhere = "WHERE id = $1"
		arg = obj.ID
	} else if obj.ExternalID != "" {
		queryWhere = "WHERE external_id = $1"
		arg = obj.ExternalID
	}

	err = QueryRow(ctx, u.db, `
SELECT 
    id, 
    nickname, 
    created_at,
    language_id,
    external_id,
	payed,
	remaining_time,
	is_redeem_used
FROM users
`+queryWhere,
		arg).Scan(
		&user.ID,
		&user.Nickname,
		&user.CreatedAt,
		&user.LanguageID,
		&user.ExternalID,
		&user.Payed,
		&user.RemainingTime,
		&user.IsRedeemUsed,
	)
	if err != nil {
		return user, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

// List
func (u *user) List(ctx context.Context, pars models.UserListPars) ([]models.User, int, error) {
	querySelect := `
   SELECT id, 
    nickname, 
	created_at,
	language_id,
	external_id,
	payed,
	remaining_time
`
	queryFrom := "FROM users"
	queryWhere := " WHERE "
	queryOffset := fmt.Sprintf(" OFFSET %d ", pars.Offset)
	queryLimit := fmt.Sprintf(" LIMIT %d ", pars.Limit)

	var args []any
	if pars.ID != nil {
		args = append(args, *pars.ID)
		queryWhere += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.Nickname != nil {
		args = append(args, *pars.Nickname)
		queryWhere += fmt.Sprintf(" nickname = $%v AND", len(args))
	}

	if pars.ExternalID != nil {
		args = append(args, *pars.ExternalID)
		queryWhere += fmt.Sprintf(" external_id = $%v AND", len(args))
	}

	queryWhere = queryWhere[:len(queryWhere)-4] // Remove the trailing " AND"

	var count int
	err := QueryRow(ctx, u.db, "select count(*) from users "+queryWhere, args...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	if pars.OnlyCount {
		return nil, count, err
	}

	rows, err := Query(ctx, u.db, querySelect+queryFrom+queryWhere+queryOffset+queryLimit, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Nickname,
			&user.CreatedAt,
			&user.LanguageID,
			&user.ExternalID,
			&user.Payed,
			&user.RemainingTime,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return users, count, nil
}

// Update
func (u *user) Update(ctx context.Context, pars models.UserGetPars, user models.UserCU) error {
	queryUpdate := `
UPDATE users
SET 
`
	var queryWhere string
	var args []any
	if pars.ID != 0 {
		queryWhere = "WHERE id = $1"
		args = append(args, pars.ID)
	} else if pars.ExternalID != "" {
		queryWhere = "WHERE external_id = $1"
		args = append(args, pars.ExternalID)
	}

	if user.Nickname != nil {
		args = append(args, *user.Nickname)
		queryUpdate += fmt.Sprintf(" nickname = $%v ,", len(args))
	}

	if user.LanguageID != nil {
		args = append(args, *user.LanguageID)
		queryUpdate += fmt.Sprintf(" language_id = $%v ,", len(args))
	}

	if user.RemainingTime != nil {
		args = append(args, *user.RemainingTime)
		queryUpdate += fmt.Sprintf(" remaining_time = remaining_time + $%v ,", len(args))
	}

	if user.Payed != nil {
		args = append(args, *user.Payed)
		queryUpdate += fmt.Sprintf(" payed = $%v ,", len(args))
	}

	if user.IsRedeemUsed != nil {
		args = append(args, *user.IsRedeemUsed)
		queryUpdate += fmt.Sprintf(" is_redeem_used = $%v ,", len(args))
	}

	queryUpdate = queryUpdate[:len(queryUpdate)-1]

	_, err := Exec(ctx, u.db, queryUpdate+queryWhere, args...)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

// Delete
func (u *user) Delete(ctx context.Context, obj models.UserGetPars) error {

	var (
		queryWhere string
		arg        any
	)
	if obj.ID != 0 {
		queryWhere = "WHERE id = $1"
		arg = obj.ID
	} else {
		queryWhere = "WHERE external_id = $1"
		arg = obj.ExternalID
	}

	_, err := Exec(ctx, u.db, "DELETE FROM users "+queryWhere, arg)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	return nil
}
