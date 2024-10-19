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
		ur   userRepo
		lang langRepo
		file fileRepo
	}

	RequiredField struct {
		ErrorCode   string
		Description string
	}
)

func New(ur userRepo, lang langRepo, file fileRepo) *Validator {
	return &Validator{
		ur:   ur,
		lang: lang,
		file: file,
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

		users, count, err := v.ur.List(ctx, models.UserListPars{Email: user.Email})
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

		users, count, err := v.ur.List(ctx, models.UserListPars{Nickname: user.Nickname})
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

	if user.LanguageID == nil && id == -1 {
		mp["language_id"] = RequiredField{
			Description: ErrMissingLanguageID.Error(),
		}
	}

	if user.LanguageID != nil {
		if *user.LanguageID == 0 {
			mp["language_id"] = RequiredField{
				Description: ErrMissingLanguageID.Error(),
			}
		}

		if _, err := v.lang.Get(ctx, *user.LanguageID); err != nil {
			mp["language_id"] = RequiredField{
				Description: ErrInvalidLanguageID.Error(),
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

func (v *Validator) ValidateLanguage(ctx context.Context, obj models.LanguageCU) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	if obj.ShortName == nil {
		mp["short_name"] = RequiredField{
			Description: ErrMissingShortName.Error(),
		}
	}
	if obj.LongName == nil {
		mp["long_name"] = RequiredField{
			Description: ErrMissingLongName.Error(),
		}
	}

	return mp, nil
}

func (v *Validator) ValidateTheme(ctx context.Context, obj models.ThemeCU) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	if obj.Level == nil {
		mp["level"] = RequiredField{
			Description: ErrMissingLevel.Error(),
		}
	}

	if obj.LanguageID == nil {
		mp["language_id"] = RequiredField{
			Description: ErrMissingLanguageID.Error(),
		}
	}

	if obj.Topic == nil {
		mp["topic"] = RequiredField{
			Description: ErrMissingTopic.Error(),
		}
	}

	if obj.Question == nil {
		mp["question"] = RequiredField{
			Description: ErrMissingQuestion.Error(),
		}
	}
	return mp, nil
}

func (v *Validator) ValidateFile(ctx context.Context, obj models.FileCU) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	if obj.Key == nil {
		mp["key"] = RequiredField{
			Description: ErrMissingFileKey.Error(),
		}
	}

	if obj.Key != nil {
		if *obj.Key == "" {
			mp["key"] = RequiredField{
				Description: ErrMissingFileKey.Error(),
			}
		}

		_, err := v.file.GetByKey(ctx, *obj.Key)
		if err == nil {
			mp["key"] = RequiredField{
				Description: ErrFileKeyIsUsed.Error(),
			}
		}
	}

	return mp, nil
}
