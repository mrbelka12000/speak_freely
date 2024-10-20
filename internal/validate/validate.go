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
		user  userRepo
		lang  langRepo
		file  fileRepo
		theme themeRepo
	}

	RequiredField struct {
		ErrorCode   string
		Description string
	}
)

func New(ur userRepo, lang langRepo, file fileRepo, theme themeRepo) *Validator {
	return &Validator{
		user:  ur,
		lang:  lang,
		file:  file,
		theme: theme,
	}
}

// ValidateUser validates create and update of users information
func (v *Validator) ValidateUser(ctx context.Context, user models.UserCU, id int64) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	// firstname check
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

	// lastname check
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

	// email
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

		users, count, err := v.user.List(ctx, models.UserListPars{Email: user.Email})
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

	// nickname check
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

		users, count, err := v.user.List(ctx, models.UserListPars{Nickname: user.Nickname})
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

	// password check
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

	// language check
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
	if obj.ShortName != nil {
		if *obj.ShortName == "" {
			mp["short_name"] = RequiredField{
				Description: ErrMissingShortName.Error(),
			}
		}
	}

	if obj.LongName == nil {
		mp["long_name"] = RequiredField{
			Description: ErrMissingLongName.Error(),
		}
	}
	if obj.LongName != nil {
		if *obj.LongName == "" {
			mp["short_name"] = RequiredField{
				Description: ErrMissingLongName.Error(),
			}
		}
	}

	return mp, nil
}

func (v *Validator) ValidateTheme(ctx context.Context, obj models.ThemeCU) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	// level check
	if obj.Level == nil {
		mp["level"] = RequiredField{
			Description: ErrMissingLevel.Error(),
		}
	}

	// language check
	if obj.LanguageID == nil {
		mp["language_id"] = RequiredField{
			Description: ErrMissingLanguageID.Error(),
		}
	}
	if obj.LanguageID != nil {
		if _, err := v.lang.Get(ctx, *obj.LanguageID); err != nil {
			mp["language_id"] = RequiredField{
				Description: ErrInvalidLanguageID.Error(),
			}
		}
	}

	// topic check
	if obj.Topic == nil {
		mp["topic"] = RequiredField{
			Description: ErrMissingTopic.Error(),
		}
	}
	if obj.Topic != nil {
		if *obj.Topic == "" {
			mp["topic"] = RequiredField{
				Description: ErrMissingTopic.Error(),
			}
		}
	}

	// question check
	if obj.Question == nil {
		mp["question"] = RequiredField{
			Description: ErrMissingQuestion.Error(),
		}
	}
	if obj.Question != nil {
		if *obj.Question == "" {
			mp["question"] = RequiredField{
				Description: ErrMissingQuestion.Error(),
			}
		}
	}

	return mp, nil
}

func (v *Validator) ValidateTranscript(ctx context.Context, obj models.TranscriptCU, id int64) (map[string]RequiredField, error) {
	mp := make(map[string]RequiredField)

	// theme check
	if obj.ThemeID == nil && id == -1 {
		mp["theme_id"] = RequiredField{
			Description: ErrMissingThemeID.Error(),
		}
	}
	if obj.ThemeID != nil {
		if *obj.ThemeID == 0 {
			mp["theme_id"] = RequiredField{
				Description: ErrMissingThemeID.Error(),
			}
		}
		if _, err := v.theme.Get(ctx, *obj.ThemeID); err != nil {
			mp["theme_id"] = RequiredField{
				Description: ErrInvalidThemeID.Error(),
			}
		}
	}

	// language check
	if obj.LanguageID == nil && id == -1 {
		mp["language_id"] = RequiredField{
			Description: ErrMissingLanguageID.Error(),
		}
	}

	if obj.LanguageID != nil {
		if *obj.LanguageID == 0 {
			mp["language_id"] = RequiredField{
				Description: ErrMissingLanguageID.Error(),
			}
		}

		if _, err := v.lang.Get(ctx, *obj.LanguageID); err != nil {
			mp["language_id"] = RequiredField{
				Description: ErrInvalidLanguageID.Error(),
			}
		}
	}

	// user check
	if obj.UserID == nil && id == -1 {
		mp["user_id"] = RequiredField{
			Description: ErrMissingUserID.Error(),
		}
	}

	if obj.UserID != nil {
		if *obj.UserID == 0 {
			mp["user_id"] = RequiredField{
				Description: ErrMissingUserID.Error(),
			}
		}

		if _, err := v.user.Get(ctx, *obj.UserID); err != nil {
			mp["user_id"] = RequiredField{
				Description: ErrInvalidUserID.Error(),
			}
		}
	}

	// file check
	if obj.FileID == nil && id == -1 {
		mp["file_id"] = RequiredField{
			Description: ErrMissingFileID.Error(),
		}
	}
	if obj.FileID != nil {
		if *obj.FileID == 0 {
			mp["file_id"] = RequiredField{
				Description: ErrMissingFileID.Error(),
			}
		}

		if _, err := v.file.Get(ctx, *obj.FileID); err != nil {
			mp["file_id"] = RequiredField{
				Description: ErrInvalidFileID.Error(),
			}
		}
	}

	// text check
	if obj.Text != nil {
		if *obj.Text == "" {
			mp["text"] = RequiredField{
				Description: ErrMissingText.Error(),
			}
		}
	}

	return mp, nil
}
