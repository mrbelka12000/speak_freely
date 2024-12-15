package validate

import (
	"context"

	"github.com/mrbelka12000/speak_freely/internal/models"
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
	if obj.TopicID == nil {
		mp["topic"] = RequiredField{
			Description: ErrMissingTopic.Error(),
		}
	}
	if obj.TopicID != nil {
		if *obj.TopicID == 0 {
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

		if _, err := v.user.Get(ctx, models.UserGetPars{ID: *obj.UserID}); err != nil {
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
				Description: ErrNoMatchingWordsInLanguage.Error(),
			}
		}
	}

	return mp, nil
}
