package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
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
			first_name,
			last_name,
			nickname,
			email,
			password,
			auth_method,
			first_language,
			created_at
			) 
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
		 	$7,
		    $8
		) RETURNING id`,
		*user.FirstName,
		*user.LastName,
		*user.Nickname,
		*user.Email,
		*user.Password,
		*user.AuthMethod,
		*user.FirstLanguage,
		user.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

// Get
func (u *user) Get(ctx context.Context, id int64) (user models.User, err error) {
	err = QueryRow(ctx, u.db, `
SELECT 
    id, 
    first_name, 
    last_name, 
    nickname, 
    email, 
    auth_method,
    created_at,
    first_language,
	confirmed
FROM users
WHERE id = $1`,
		id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Nickname,
		&user.Email,
		&user.AuthMethod,
		&user.CreatedAt,
		&user.FirstLanguage,
		&user.Confirmed,
	)
	if err != nil {
		return user, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

// List
func (u *user) List(ctx context.Context, pars models.UserPars) ([]models.User, int, error) {
	querySelect := `
   SELECT id, 
    first_name, 
    last_name, 
    nickname, 
    email, 
    auth_method,
	created_at,
	first_language,
	confirmed
`
	queryFrom := "FROM users"
	queryWhere := " WHERE "
	var args []any
	if pars.ID != nil {
		args = append(args, *pars.ID)
		queryWhere += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.FirstName != nil {
		args = append(args, *pars.FirstName)
		queryWhere += fmt.Sprintf(" first_name = $%v AND", len(args))
	}

	if pars.LastName != nil {
		args = append(args, *pars.LastName)
		queryWhere += fmt.Sprintf(" last_name = $%v AND", len(args))
	}

	if pars.Nickname != nil {
		args = append(args, *pars.Nickname)
		queryWhere += fmt.Sprintf(" nickname = $%v AND", len(args))
	}

	if pars.Email != nil {
		args = append(args, *pars.Email)
		queryWhere += fmt.Sprintf(" email = $%v AND", len(args))
	}

	if pars.Confirmed != nil {
		args = append(args, *pars.Confirmed)
		queryWhere += fmt.Sprintf(" confirmed = $%v AND", len(args))

	}

	queryWhere = queryWhere[:len(queryWhere)-4] // Remove the trailing " AND"

	if pars.Count {
		var count int

		err := QueryRow(ctx, u.db, "select count(*) from users "+queryWhere, args...).Scan(&count)
		if err != nil {
			return nil, 0, fmt.Errorf("count users: %w", err)
		}

		return nil, count, err
	}

	rows, err := Query(ctx, u.db, querySelect+queryFrom+queryWhere, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.Email,
			&user.AuthMethod,
			&user.CreatedAt,
			&user.FirstLanguage,
			&user.Confirmed,
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

	return users, len(users), nil
}

// Update
func (u *user) Update(ctx context.Context, id int64, user models.UserCU) error {
	queryUpdate := `
UPDATE users
SET 
`

	queryWhere := "WHERE id = $1"

	var args []any
	args = append(args, id)

	if user.FirstName != nil {
		args = append(args, *user.FirstName)
		queryUpdate += fmt.Sprintf(" first_name = $%v ,", len(args))
	}

	if user.LastName != nil {
		args = append(args, *user.LastName)
		queryUpdate += fmt.Sprintf(" last_name = $%v ,", len(args))
	}

	if user.Nickname != nil {
		args = append(args, *user.Nickname)
		queryUpdate += fmt.Sprintf(" nickname = $%v ,", len(args))
	}

	if user.Email != nil {
		args = append(args, *user.Email)
		queryUpdate += fmt.Sprintf(" email = $%v ,", len(args))
	}

	if user.Confirmed {
		args = append(args, user.Confirmed)
		queryUpdate += fmt.Sprintf(" confirmed = $%v ,", len(args))
	}

	queryUpdate = queryUpdate[:len(queryUpdate)-1]

	_, err := Exec(ctx, u.db, queryUpdate+queryWhere, args...)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

// Delete
func (u *user) Delete(ctx context.Context, id int64) error {

	_, err := Exec(ctx, u.db, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	return nil
}

// FindByLogin
func (u *user) FindByLogin(ctx context.Context, login string) (out models.User, err error) {
	err = QueryRow(ctx, u.db, `
SELECT 
    id, 
    first_name, 
    last_name, 
    nickname, 
    email, 
    auth_method,
    created_at,
    password
FROM users
WHERE nickname = $1 OR email = $1 AND auth_method = 1`,
		login).Scan(
		&out.ID,
		&out.FirstName,
		&out.LastName,
		&out.Nickname,
		&out.Email,
		&out.AuthMethod,
		&out.CreatedAt,
		&out.Password,
	)
	if err != nil {
		return out, fmt.Errorf("get user: %w", err)
	}

	return out, nil
}
