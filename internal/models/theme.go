package models

type (
	// Theme
	Theme struct {
		ID         int64  `json:"id,omitempty"`
		LanguageID int64  `json:"language_id,omitempty"`
		Topic      string `json:"topic,omitempty"`
		Question   string `json:"question,omitempty"`
		Level      string `json:"level,omitempty"`

		Language *Language `json:"language,omitempty"`
	}

	ThemeCU struct {
		LanguageID *int64  `json:"language_id,omitempty"`
		Topic      *string `json:"topic,omitempty"`
		Question   *string `json:"question,omitempty"`
		Level      *string `json:"level,omitempty"`
	}

	ThemeListPars struct {
		ID         *int64  `json:"id,omitempty" schema:"id"`
		LanguageID *int64  `json:"language_id,omitempty" schema:"language_id"`
		Level      *string `json:"level,omitempty" schema:"level"`

		OnlyCount        bool `json:"only_count,omitempty" schema:"only_count"`
		PaginationParams `json:"pagination_params"`
	}
)
