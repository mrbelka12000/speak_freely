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
		Text       *string
		LanguageID *int64
		UserID     *int64
		FileID     *int64
		ThemeID    *int64
		Accuracy   *float64
	}

	TranscriptListPars struct {
		ID         *int64
		LanguageID *int64
		UserID     *int64
		ThemeID    *int64

		OnlyCount bool

		PaginationParams
	}
)
