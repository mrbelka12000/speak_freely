package validate

import (
	"context"
	"fmt"
	"net/mail"
	"unicode"
	"unicode/utf8"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

type (
	Validator struct {
		ur userRepo
	}

	RequiredField struct {
		ErrorCode   string
		Description string
	}
)

func New(ur userRepo) *Validator {
	return &Validator{
		ur: ur,
	}
}

// ValidateUser validates create and update of users information
func (v *Validator) ValidateUser(ctx context.Context, user models.UserCU, id int64) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	if user.FirstName == nil && id == -1 {
		mp["first_name"] = RequiredField{
			Description: ErrMissingFirstName.Error(),
		}
	}
	if user.FirstName != nil {
		if *user.FirstName == "" {
			mp["first_name"] = RequiredField{
				Description: ErrMissingFirstName.Error(),
			}
		}

		if utf8.RuneCountInString(*user.FirstName) > 50 {
			mp["first_name"] = RequiredField{
				Description: ErrFirstNameTooLong.Error(),
			}
		}
	}

	if user.LastName == nil && id == -1 {
		mp["last_name"] = RequiredField{
			Description: ErrMissingLastName.Error(),
		}
	}
	if user.LastName != nil {
		if *user.LastName == "" {
			mp["last_name"] = RequiredField{
				Description: ErrMissingLastName.Error(),
			}
		}

		if utf8.RuneCountInString(*user.LastName) > 50 {
			mp["last_name"] = RequiredField{
				Description: ErrLastNameTooLong.Error(),
			}
		}
	}

	if user.Email == nil && id == -1 {
		mp["email"] = RequiredField{
			Description: ErrMissingEmail.Error(),
		}
	}
	if user.Email != nil {
		if *user.Email == "" {
			mp["email"] = RequiredField{
				Description: ErrMissingEmail.Error(),
			}
		}

		if utf8.RuneCountInString(*user.Email) > 100 {
			mp["email"] = RequiredField{
				Description: ErrEmailTooLong.Error(),
			}
		}

		_, err := mail.ParseAddress(*user.Email)
		if err != nil {
			mp["email"] = RequiredField{
				Description: ErrInvalidEmail.Error(),
			}
		}

		users, count, err := v.ur.List(ctx, models.UserPars{Email: user.Email})
		if err != nil {
			return nil, fmt.Errorf("get user list: %w", err)
		}

		if id == -1 {
			if count > 0 {
				mp["email"] = RequiredField{
					Description: ErrEmailIsUsed.Error(),
				}
			}
		} else {
			for _, user := range users {
				if user.ID != id {
					mp["email"] = RequiredField{
						Description: ErrEmailIsUsed.Error(),
					}
				}
			}
		}
	}

	if user.Nickname == nil && id == -1 {
		mp["nickname"] = RequiredField{
			Description: ErrMissingNickname.Error(),
		}
	}
	if user.Nickname != nil {
		if *user.Nickname == "" {
			mp["nickname"] = RequiredField{
				Description: ErrMissingNickname.Error(),
			}
		}
		if utf8.RuneCountInString(*user.Nickname) > 50 {
			mp["nickname"] = RequiredField{
				Description: ErrNicknameTooLong.Error(),
			}
		}

		users, count, err := v.ur.List(ctx, models.UserPars{Nickname: user.Nickname})
		if err != nil {
			return nil, fmt.Errorf("get user list: %w", err)
		}
		if id == -1 {
			if count > 0 {
				mp["nickname"] = RequiredField{
					Description: ErrNicknameIsUsed.Error(),
				}
			}
		} else {
			for _, user := range users {
				if user.ID != id {
					mp["nickname"] = RequiredField{
						Description: ErrNicknameIsUsed.Error(),
					}
				}
			}
		}
	}

	if user.Password == nil && id == -1 {
		mp["password"] = RequiredField{
			Description: ErrMissingPassword.Error(),
		}
	}
	if user.Password != nil {
		if *user.Password == "" {
			mp["password"] = RequiredField{
				Description: ErrMissingPassword.Error(),
			}
		}

		if err := validatePassword(*user.Password); err != nil {
			mp["password"] = RequiredField{
				Description: err.Error(),
			}
		}
	}

	if user.FirstLanguage == nil && id == -1 {
		mp["first_language"] = RequiredField{
			Description: ErrMissingFirstLanguage.Error(),
		}
	}

	if user.FirstLanguage != nil {
		if *user.FirstLanguage == "" {
			mp["first_language"] = RequiredField{
				Description: ErrMissingFirstLanguage.Error(),
			}
		}
		if utf8.RuneCountInString(*user.FirstLanguage) > 3 {
			mp["first_language"] = RequiredField{
				Description: ErrFirstLanguageTooLong.Error(),
			}
		}
	}

	return mp, nil
}

func validatePassword(password string) (err error) {
	var (
		anyNum    bool
		anyLetter bool
	)

	for _, c := range password {
		if unicode.IsSymbol(c) || unicode.IsPunct(c) || unicode.IsSpace(c) {
			return ErrInvalidPassword
		}
		if unicode.IsLetter(c) {
			anyLetter = true
		}
		if unicode.IsNumber(c) {
			anyNum = true
		}
	}

	if anyLetter && anyNum {
		return nil
	}

	if !anyLetter {
		return ErrNeedLetter
	}

	if !anyNum {
		return ErrNeedNumber
	}

	return nil
}
