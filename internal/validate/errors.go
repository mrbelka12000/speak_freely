package validate

import "errors"

var (
	// General errors
	ErrTooLong    = errors.New("too long")
	ErrTooShort   = errors.New("too short")
	ErrNeedLetter = errors.New("need letter in password")
	ErrNeedNumber = errors.New("need number in password")

	// User validation
	ErrMissingFirstName  = errors.New("first name not provided")
	ErrMissingLastName   = errors.New("last name not provided")
	ErrMissingNickname   = errors.New("nickname not provided")
	ErrMissingEmail      = errors.New("email not provided")
	ErrMissingPassword   = errors.New("password not provided")
	ErrMissingLanguageID = errors.New("language id not provided")
	ErrMissingAuthMethod = errors.New("auth method not provided")
	ErrInvalidFirstName  = errors.New("invalid first name")
	ErrFirstNameTooLong  = errors.New("first name too long")
	ErrLastNameTooLong   = errors.New("last name too long")
	ErrNicknameTooLong   = errors.New("nickname too long")
	ErrEmailTooLong      = errors.New("email too long")
	ErrInvalidLastName   = errors.New("invalid last name")
	ErrInvalidNickname   = errors.New("invalid nickname")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidLanguageID = errors.New("invalid language id")
	ErrInvalidAuthMethod = errors.New("invalid auth method")
	ErrEmailIsUsed       = errors.New("email is used")
	ErrNicknameIsUsed    = errors.New("nickname is used")

	// Language validation
	ErrMissingShortName = errors.New("short name not provided")
	ErrMissingLongName  = errors.New("long name not provided")

	// Theme validation
	ErrMissingLevel    = errors.New("level not provided")
	ErrMissingTopic    = errors.New("topic not provided")
	ErrMissingQuestion = errors.New("question not provided")

	// Transcript validation
	ErrMissingThemeID = errors.New("theme id not provided")
	ErrInvalidThemeID = errors.New("invalid theme id")
	ErrMissingUserID  = errors.New("user id not provided")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrMissingFileID  = errors.New("file id not provided")
	ErrInvalidFileID  = errors.New("file id not provided")
	ErrMissingText    = errors.New("text not provided")
)
