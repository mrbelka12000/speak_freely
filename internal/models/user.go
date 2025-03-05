package models

type (
	// User general user information
	User struct {
		ID            int64  `json:"id"`
		Nickname      string `json:"nickname"`
		CreatedAt     int64  `json:"created_at"`
		LanguageID    int64  `json:"language_id"`
		ExternalID    string `json:"external_id"`
		Payed         bool   `json:"payed"`
		RemainingTime int64  `json:"remaining_time"`
		IsRedeemUsed  bool   `json:"used_redeem"`
		IsStarted     bool   `json:"started"`

		Language *Language `json:"first_language,omitempty"`
	}

	UserGetPars struct {
		ID         int64  `json:"id"`
		ExternalID string `json:"external_id"`
	}

	// UserCU object to create/update user information
	UserCU struct {
		Nickname      *string `json:"nickname"`
		LanguageID    *int64  `json:"language_id"`
		ExternalID    *string `json:"external_id"`
		RemainingTime *int64  `json:"remaining_time"`
		Payed         *bool   `json:"payed"`
		IsRedeemUsed  *bool   `json:"used_redeem"`
		IsStarted     *bool   `json:"started"`
		CreatedAt     int64   `json:"created_at"`
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
