package usecase

import (
	"context"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
	"github.com/mrbelka12000/speak_freely/internal/validate"
	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

func (uc *UseCase) ThemeGet(ctx context.Context, id int64) (models.Theme, error) {
	theme, err := uc.srv.Theme.Get(ctx, id)
	if err != nil {
		return models.Theme{}, err
	}

	lang, err := uc.srv.Language.Get(ctx, theme.LanguageID)
	if err != nil {
		return models.Theme{}, err
	}
	theme.Language = &lang

	return theme, nil
}

func (uc *UseCase) ThemeBuild(ctx context.Context, obj models.ThemeCU) (int64, map[string]validate.RequiredField, error) {
	missed, err := uc.validator.ValidateTheme(ctx, obj)
	if err != nil {
		return 0, nil, fmt.Errorf("validate theme %w", err)
	}
	if len(missed) > 0 {
		return 0, missed, nil
	}

	id, err := uc.srv.Theme.Create(ctx, obj)
	if err != nil {
		return 0, nil, fmt.Errorf("create theme %w", err)
	}

	return id, nil, nil
}

func (uc *UseCase) ThemesGenerateWithAI(ctx context.Context, level string) error {
	topics, err := uc.gen.GenerateTopics(ctx, ai.GenerateTopicsRequest{
		Level: level,
	})
	if err != nil {
		return fmt.Errorf("failed to generate topics: %w", err)
	}

	ctx, err = uc.tx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer uc.tx.Rollback(ctx)

	for _, topic := range topics {
		lang, err := uc.srv.Language.GetByShortName(ctx, topic.Lang)
		if err != nil {
			return fmt.Errorf("failed to get language by name: %w", err)
		}

		_, err = uc.srv.Theme.Create(ctx, models.ThemeCU{
			LanguageID: pointer.Of(lang.ID),
			Topic:      pointer.Of(topic.Topic),
			Question:   pointer.Of(topic.Question),
			Level:      pointer.Of(level),
		})
		if err != nil {
			return fmt.Errorf("failed to create theme: %w", err)
		}
	}

	err = uc.tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (uc *UseCase) ThemeList(ctx context.Context, pars models.ThemeListPars) ([]models.Theme, int, error) {
	return uc.srv.Theme.List(ctx, pars)
}
