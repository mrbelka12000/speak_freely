package validate

import "errors"

type CustomErrors error

var (
	// General errors
	ErrTooLong    CustomErrors = errors.New("too long")
	ErrTooShort   CustomErrors = errors.New("too short")
	ErrNeedLetter CustomErrors = errors.New("need letter in password")
	ErrNeedNumber CustomErrors = errors.New("need number in password")

	// User validation
	ErrMissingFirstName  CustomErrors = errors.New("first name not provided")
	ErrMissingLastName   CustomErrors = errors.New("last name not provided")
	ErrMissingNickname   CustomErrors = errors.New("nickname not provided")
	ErrMissingEmail      CustomErrors = errors.New("email not provided")
	ErrMissingPassword   CustomErrors = errors.New("password not provided")
	ErrMissingAuthMethod CustomErrors = errors.New("auth method not provided")
	ErrInvalidFirstName  CustomErrors = errors.New("invalid first name")
	ErrFirstNameTooLong  CustomErrors = errors.New("first name too long")
	ErrLastNameTooLong   CustomErrors = errors.New("last name too long")
	ErrNicknameTooLong   CustomErrors = errors.New("nickname too long")
	ErrEmailTooLong      CustomErrors = errors.New("email too long")
	ErrInvalidLastName   CustomErrors = errors.New("invalid last name")
	ErrInvalidNickname   CustomErrors = errors.New("invalid nickname")
	ErrInvalidEmail      CustomErrors = errors.New("invalid email")
	ErrInvalidPassword   CustomErrors = errors.New("invalid password")
	ErrInvalidAuthMethod CustomErrors = errors.New("invalid auth method")
	ErrEmailIsUsed       CustomErrors = errors.New("email is used")
	ErrNicknameIsUsed    CustomErrors = errors.New("nickname is used")
)
