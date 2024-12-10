package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type theme struct {
	db *sql.DB
}

func newTheme(db *sql.DB) *theme {
	return &theme{db: db}
}

func (t *theme) Create(ctx context.Context, obj models.ThemeCU) (id int64, err error) {
	err = QueryRow(ctx, t.db, `
	INSERT INTO themes (
        level, 
		topic, 
        question,
		language_id) 
	VALUES (
        $1,
        $2,
        $3,
        $4
	) RETURNING id`,
		*obj.Level,
		*obj.Topic,
		*obj.Question,
		*obj.LanguageID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create theme: %w", err)
	}

	return id, nil
}

func (t *theme) Get(ctx context.Context, id int64) (theme models.Theme, err error) {
	err = QueryRow(ctx, t.db, `
SELECT 
	id,
	level,
	topic,
	question,
	language_id
FROM themes
WHERE id = $1`,
		id).Scan(
		&theme.ID,
		&theme.Level,
		&theme.Topic,
		&theme.Question,
		&theme.LanguageID,
	)
	if err != nil {
		return theme, fmt.Errorf("get theme: %w", err)
	}

	return theme, nil
}

func (t *theme) List(ctx context.Context, pars models.ThemeListPars) ([]models.Theme, int, error) {
	querySelect := `
   SELECT DISTINCT ON (question)
	id,
	level,
	topic,
	question,
	language_id
`
	queryFrom := "FROM themes"
	queryWhere := " WHERE "
	queryOffset := fmt.Sprintf(" OFFSET %d ", pars.Offset)
	queryLimit := fmt.Sprintf(" LIMIT %d ", pars.Limit)
	var queryOrderBy string
	var args []any
	if pars.ID != nil {
		args = append(args, *pars.ID)
		queryWhere += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.Level != nil {
		args = append(args, *pars.Level)
		queryWhere += fmt.Sprintf(" level = $%v AND", len(args))
	}

	if pars.LanguageID != nil {
		args = append(args, *pars.LanguageID)
		queryWhere += fmt.Sprintf(" language_id = $%v AND", len(args))
	}

	if pars.Topic != nil {
		args = append(args, *pars.Topic)
		queryWhere += fmt.Sprintf(" topic = $%v AND", len(args))
	}

	queryWhere = queryWhere[:len(queryWhere)-4] // Remove the trailing " AND"

	var count int
	err := QueryRow(ctx, t.db, "select count(*) from themes "+queryWhere, args...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if pars.OnlyCount {
		return nil, count, nil
	}

	if pars.Random {
		queryOrderBy = " ORDER BY random()"
	}

	fmt.Println(querySelect + queryFrom + queryWhere + queryOrderBy + queryOffset + queryLimit)
	rows, err := Query(ctx, t.db, querySelect+queryFrom+queryWhere+queryOrderBy+queryOffset+queryLimit, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var themes []models.Theme

	for rows.Next() {
		var theme models.Theme
		err := rows.Scan(
			&theme.ID,
			&theme.Level,
			&theme.Topic,
			&theme.Question,
			&theme.LanguageID,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan theme: %w", err)
		}
		themes = append(themes, theme)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return themes, count, nil
}
