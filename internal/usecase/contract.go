package usecase

import (
	"context"
	"io"
	"time"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
)

type (
	txBuilder interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}

	cache interface {
		Set(key string, value interface{}, dur time.Duration) error
		Get(key string) (string, bool)
		GetInt64(key string) (int64, bool)
		Delete(key string)
	}

	generator interface {
		GenerateTheme(ctx context.Context, req ai.GenerateThemeRequest) (obj ai.GenerateThemeResponse, err error)
		GetSuggestions(ctx context.Context, req ai.SuggestionRequest) (obj ai.SuggestionResponse, err error)
		Dialog(ctx context.Context, req ai.DialogRequest) (obj ai.DialogResponse, err error)
	}

	storage interface {
		UploadFile(ctx context.Context, file io.Reader, objectName, contentType string, fileSize int64) (string, error)
	}

	transcriber interface {
		GetTextFromFile(ctx context.Context, file io.Reader, lang string) (string, error)
		GetTextFromURL(ctx context.Context, url, languageCode string) (string, error)
	}
)
