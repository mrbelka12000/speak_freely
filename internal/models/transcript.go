package models

type (
	Transcript struct {
		ID         int64
		Text       *string
		Accuracy   *float64 // will be filled after calculations
		LanguageID int64
		UserID     int64
		FileID     int64
		ThemeID    int64
	}

	TranscriptCU struct {
		Text       *string  `json:"text,omitempty"`
		LanguageID *int64   `json:"language_id,omitempty"`
		UserID     *int64   `json:"user_id,omitempty"`
		FileID     *int64   `json:"file_id,omitempty"`
		ThemeID    *int64   `json:"theme_id,omitempty"`
		Accuracy   *float64 `json:"accuracy,omitempty"`
	}

	TranscriptListPars struct {
		ID         *int64 `json:"id,omitempty" schema:"id"`
		LanguageID *int64 `json:"language_id,omitempty" schema:"language_id"`
		UserID     *int64 `json:"user_id,omitempty" schema:"user_id"`
		ThemeID    *int64 `json:"theme_id,omitempty" schema:"theme_id"`

		OnlyCount bool `json:"only_count,omitempty" schema:"only_count"`

		PaginationParams `json:"pagination_params"`
	}
)
