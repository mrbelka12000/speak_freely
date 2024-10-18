package models

type (
	// User general user information
	User struct {
		ID         int64  `json:"id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Nickname   string `json:"nickname"`
		Email      string `json:"email"`
		Password   string `json:"password,omitempty"`
		AuthMethod int    `json:"auth_method"`
		CreatedAt  int64  `json:"created_at"`
		Confirmed  bool   `json:"confirmed"`
		LanguageID int64  `json:"language_id"`

		Language *Language `json:"first_language,omitempty"`
	}

	// UserLogin for login to the website
	UserLogin struct {
		Login    string `json:"login"` // nickname/email
		Password string `json:"password"`
	}

	// UserCU object to create/update user information
	UserCU struct {
		FirstName  *string `json:"first_name"`
		LastName   *string `json:"last_name"`
		Nickname   *string `json:"nickname"`
		Email      *string `json:"email"`
		Password   *string `json:"password,omitempty"`
		AuthMethod *int    `json:"auth_method"`
		LanguageID *int64  `json:"language_id"`
		CreatedAt  int64   `json:"created_at"`
		Confirmed  bool    `json:"confirmed"`
	}

	// UserListPars for list users
	UserListPars struct {
		ID        *int64
		FirstName *string
		LastName  *string
		Nickname  *string
		Email     *string
		Confirmed *bool
		Count     bool
	}

	// UserGetPars
	UserGetPars struct {
		ID int64
	}
)
