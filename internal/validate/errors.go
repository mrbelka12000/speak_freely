package validate

import "errors"

var (
	// General errors
	ErrTooLong    = errors.New("too long")
	ErrTooShort   = errors.New("too short")
	ErrNeedLetter = errors.New("need letter in password")
	ErrNeedNumber = errors.New("need number in password")

	// User validation
	ErrMissingFirstName     = errors.New("first name not provided")
	ErrMissingLastName      = errors.New("last name not provided")
	ErrMissingNickname      = errors.New("nickname not provided")
	ErrMissingEmail         = errors.New("email not provided")
	ErrMissingPassword      = errors.New("password not provided")
	ErrMissingFirstLanguage = errors.New("first language not provided")
	ErrMissingAuthMethod    = errors.New("auth method not provided")
	ErrInvalidFirstName     = errors.New("invalid first name")
	ErrFirstNameTooLong     = errors.New("first name too long")
	ErrLastNameTooLong      = errors.New("last name too long")
	ErrNicknameTooLong      = errors.New("nickname too long")
	ErrEmailTooLong         = errors.New("email too long")
	ErrFirstLanguageTooLong = errors.New("first language too long")
	ErrInvalidLastName      = errors.New("invalid last name")
	ErrInvalidNickname      = errors.New("invalid nickname")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidAuthMethod    = errors.New("invalid auth method")
	ErrEmailIsUsed          = errors.New("email is used")
	ErrNicknameIsUsed       = errors.New("nickname is used")
)
