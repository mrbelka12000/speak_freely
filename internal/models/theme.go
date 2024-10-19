package models

type (
	// Theme
	Theme struct {
		ID         int64  `json:"id,omitempty"`
		LanguageID int64  `json:"language_id,omitempty"`
		Topic      string `json:"topic,omitempty"`
		Question   string `json:"question,omitempty"`
		Level      int    `json:"level,omitempty"`

		Language *Language `json:"language,omitempty"`
	}

	ThemeListPars struct {
		LanguageID *int64
		Level      *int

		PaginationParams
	}
)
