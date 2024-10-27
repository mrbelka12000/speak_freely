package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

type transcript struct {
	db *sql.DB
}

func newTranscript(db *sql.DB) *transcript {
	return &transcript{
		db: db,
	}
}

// Create
func (t *transcript) Create(ctx context.Context, obj models.TranscriptCU) (id int64, err error) {
	err = QueryRow(ctx, t.db, `
	INSERT INTO transcripts(
		text,
		language_id,
		user_id,
		file_id,
		theme_id,
	    accuracy,
	    suggestion
	) VALUES (
	    $1,
		$2,
		$3,
		$4,
	    $5,
	    $6,
	    $7
	) RETURNING id
`,
		pointer.Value(obj.Text),
		pointer.Value(obj.LanguageID),
		pointer.Value(obj.UserID),
		pointer.Value(obj.FileID),
		pointer.Value(obj.ThemeID),
		pointer.Value(obj.Accuracy),
		obj.Suggestion,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("transcript create: %w", err)
	}

	return id, nil
}

// Get
func (t *transcript) Get(ctx context.Context, id int64) (obj models.Transcript, err error) {
	err = QueryRow(ctx, t.db, `
	SELECT
	id,
	text,
	accuracy,
	language_id,
	user_id,
	file_id,
	theme_id
FROM transcripts
WHERE id = $1
`, id).Scan(
		&obj.ID,
		&obj.Text,
		&obj.Accuracy,
		&obj.LanguageID,
		&obj.UserID,
		&obj.FileID,
		&obj.ThemeID,
	)
	if err != nil {
		return obj, fmt.Errorf("get transcript: %w", err)
	}

	return obj, nil
}

// List
func (t *transcript) List(ctx context.Context, pars models.TranscriptListPars) ([]models.Transcript, int, error) {
	querySelect := `
	SELECT
	id,
	text,
	accuracy,
	language_id,
	user_id,
	file_id,
	theme_id
`
	queryFrom := "FROM transcripts"
	queryWhere := " WHERE "
	queryOffset := fmt.Sprintf(" OFFSET %d ", pars.Offset)
	queryLimit := fmt.Sprintf(" LIMIT %d ", pars.Limit)

	var args []any
	if pars.ID != nil {
		args = append(args, *pars.ID)
		queryWhere += fmt.Sprintf(" id = $%v AND", len(args))
	}

	if pars.UserID != nil {
		args = append(args, *pars.UserID)
		queryWhere += fmt.Sprintf(" user_id = $%v AND", len(args))
	}

	if pars.ThemeID != nil {
		args = append(args, *pars.ThemeID)
		queryWhere += fmt.Sprintf(" theme_id = $%v AND", len(args))
	}

	if pars.LanguageID != nil {
		args = append(args, *pars.LanguageID)
		queryWhere += fmt.Sprintf(" language_id = $%v AND", len(args))
	}

	queryWhere = queryWhere[:len(queryWhere)-4] // Remove the trailing " AND"

	var count int

	err := QueryRow(ctx, t.db, "select count(*) from transcripts "+queryWhere, args...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	if pars.OnlyCount {
		return nil, count, err
	}

	rows, err := Query(ctx, t.db, querySelect+queryFrom+queryWhere+queryOffset+queryLimit, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var result []models.Transcript
	for rows.Next() {
		var obj models.Transcript
		err = rows.Scan(
			&obj.ID,
			&obj.Text,
			&obj.Accuracy,
			&obj.LanguageID,
			&obj.UserID,
			&obj.FileID,
			&obj.ThemeID,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transcript: %w", err)
		}
		result = append(result, obj)
	}

	return result, count, nil
}

// Update
func (t *transcript) Update(ctx context.Context, id int64, obj models.TranscriptCU) error {
	queryUpdate := `
UPDATE transcripts
SET 
`

	queryWhere := "WHERE id = $1"
	var args []any
	args = append(args, id)

	if obj.UserID != nil {
		args = append(args, *obj.UserID)
		queryWhere += fmt.Sprintf(" user_id = $%v ", len(args))
	}

	if obj.FileID != nil {
		args = append(args, *obj.FileID)
		queryWhere += fmt.Sprintf(" file_id = $%v ", len(args))
	}

	if obj.ThemeID != nil {
		args = append(args, *obj.ThemeID)
		queryWhere += fmt.Sprintf(" theme_id = $%v ", len(args))
	}

	if obj.LanguageID != nil {
		args = append(args, *obj.LanguageID)
		queryWhere += fmt.Sprintf(" language_id = $%v ", len(args))
	}

	if obj.Accuracy != nil {
		args = append(args, *obj.Accuracy)
		queryWhere += fmt.Sprintf(" accuracy = $%v ", len(args))
	}
	queryUpdate = queryUpdate[:len(queryUpdate)-1]

	_, err := Exec(ctx, t.db, queryUpdate+queryWhere, args...)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

// Delete
func (t *transcript) Delete(ctx context.Context, id int64) error {

	_, err := Exec(ctx, t.db, "DELETE FROM public.transcripts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	return nil
}
