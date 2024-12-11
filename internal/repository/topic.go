package repository

import (
	"context"
	"database/sql"

	"github.com/mrbelka12000/speak_freely/internal/models"
)

type topic struct {
	db *sql.DB
}

func newTopic(db *sql.DB) *topic {
	return &topic{
		db: db,
	}
}

func (t *topic) Get(ctx context.Context, id int64) (obj models.Topic, err error) {
	err = t.db.QueryRowContext(ctx, `
SELECT id, name, language_id FROM topics WHERE id = $1
`, id).Scan(&obj.ID, &obj.Name, &obj.LanguageID)
	return obj, err
}

func (t *topic) Create(ctx context.Context, obj models.Topic) error {
	_, err := t.db.ExecContext(ctx, `
INSERT INTO topics
(name, language_id)
VALUES 
($1, $2)
`, obj.Name, obj.LanguageID)
	return err
}

func (t *topic) List(ctx context.Context, languageID int64) ([]models.Topic, error) {
	queryWhere := " WHERE  1 = 1  AND"
	var args []any
	if languageID > 0 {
		queryWhere += " language_id = $1 AND"
		args = append(args, languageID)
	}

	queryWhere = queryWhere[:len(queryWhere)-4] // Remove the trailing " AND"

	rows, err := t.db.QueryContext(ctx, `
SELECT id, name, language_id
FROM topics
`+queryWhere, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		topic := models.Topic{}
		err = rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.LanguageID,
		)
		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}
