package usecase

import (
	"context"
	"fmt"
	"math/rand"

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

	topic, err := uc.srv.Topic.Get(ctx, theme.TopicID)
	if err != nil {
		return models.Theme{}, err
	}
	theme.Topic = &topic

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

	languages, _, err := uc.srv.Language.List(ctx)
	if err != nil {
		return fmt.Errorf("get language list %w", err)
	}

	for _, language := range languages {
		topics, err := uc.srv.Topic.List(ctx, language.ID)
		if err != nil {
			return fmt.Errorf("get topic list %w", err)
		}

		randTopic := topics[rand.Intn(len(topics))]

		theme, err := uc.gen.GenerateTheme(ctx, ai.GenerateThemeRequest{
			Level:    level,
			Language: language.ShortName,
			Topic:    randTopic.Name,
		})
		if err != nil {
			return fmt.Errorf("generate theme %w", err)
		}

		_, err = uc.srv.Theme.Create(ctx, models.ThemeCU{
			LanguageID: pointer.Of(language.ID),
			TopicID:    pointer.Of(randTopic.ID),
			Question:   pointer.Of(theme.Question),
			Level:      pointer.Of(level),
		})
		if err != nil {
			return fmt.Errorf("failed to create theme: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) ThemeList(ctx context.Context, pars models.ThemeListPars) ([]models.Theme, int, error) {
	return uc.srv.Theme.List(ctx, pars)
}
