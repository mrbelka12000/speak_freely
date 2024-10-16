package models

type (
	User struct {
		ID         int64  `json:"id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Nickname   string `json:"nickname"`
		Email      string `json:"email"`
		Password   string `json:"password,omitempty"`
		AuthMethod int    `json:"auth_method"`
		CreatedAt  int64  `json:"created_at"`
	}

	UserLogin struct {
		Login    string `json:"login"` // nickname/email
		Password string `json:"password"`
	}

	UserCU struct {
		FirstName  *string `json:"first_name"`
		LastName   *string `json:"last_name"`
		Nickname   *string `json:"nickname"`
		Email      *string `json:"email"`
		Password   *string `json:"password,omitempty"`
		AuthMethod *int    `json:"auth_method"`
		CreatedAt  int64   `json:"created_at"`
	}

	UserPars struct {
		ID        *int64
		FirstName *string
		LastName  *string
		Nickname  *string
		Email     *string
	}
)
