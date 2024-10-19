package usecase

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/validate"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

// SaveFile
func (uc *UseCase) SaveFile(
	ctx context.Context,
	b *bytes.Buffer,
	objectName,
	contentType string,
	fileSize int64,
) (int64, map[string]validate.RequiredField, error) {
	objectName = fmt.Sprintf("%d-%s", time.Now().UnixMilli(), objectName)

	fileKey, err := uc.storage.UploadFile(ctx, b, objectName, contentType, fileSize)
	if err != nil {
		return 0, nil, fmt.Errorf("upload file to storage: %w", err)
	}

	obj := models.FileCU{
		Key: pointer.Of(fileKey),
	}

	missed, err := uc.validator.ValidateFile(ctx, obj)
	if err != nil {
		return 0, nil, fmt.Errorf("validate file: %w", err)
	}
	if len(missed) > 0 {
		return 0, missed, nil
	}

	id, err := uc.srv.File.Create(ctx, obj)
	if err != nil {
		return 0, nil, fmt.Errorf("create file: %w", err)
	}

	return id, nil, nil
}
