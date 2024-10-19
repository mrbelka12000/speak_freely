package models

type (
	Language struct {
		ID        int64  `json:"id,omitempty"`
		LongName  string `json:"long_name,omitempty"`
		ShortName string `json:"short_name,omitempty"`
	}

	LanguageCU struct {
		LongName  *string `json:"long_name"`
		ShortName *string `json:"short_name"`
	}
)
