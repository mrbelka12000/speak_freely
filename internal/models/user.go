package models

type (
	// User general user information
	User struct {
		ID         int64  `json:"id"`
		Nickname   string `json:"nickname"`
		Password   string `json:"password,omitempty"`
		CreatedAt  int64  `json:"created_at"`
		LanguageID int64  `json:"language_id"`
		ExternalID string `json:"external_id"`

		Language *Language `json:"first_language,omitempty"`
	}

	UserGetPars struct {
		ID         int64  `json:"id"`
		ExternalID string `json:"external_id"`
	}

	// UserCU object to create/update user information
	UserCU struct {
		Nickname   *string `json:"nickname"`
		LanguageID *int64  `json:"language_id"`
		CreatedAt  int64   `json:"created_at"`
		ExternalID *string `json:"external_id"`
	}

	// UserListPars for list users
	UserListPars struct {
		ID         *int64  `json:"id,omitempty" schema:"id"`
		Nickname   *string `json:"nickname,omitempty" schema:"nickname"`
		OnlyCount  bool    `json:"only_count,omitempty" schema:"only_count"`
		ExternalID *string `json:"external_id,omitempty" schema:"external_id"`

		PaginationParams
	}
)
