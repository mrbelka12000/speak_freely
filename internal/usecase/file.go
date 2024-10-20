package usecase

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

// SaveFile
func (uc *UseCase) SaveFile(
	ctx context.Context,
	b *bytes.Buffer,
	objectName,
	contentType string,
	fileSize int64,
) (int64, error) {
	objectName = fmt.Sprintf("%d-%s", time.Now().UnixMilli(), objectName)

	fileKey, err := uc.storage.UploadFile(ctx, b, objectName, contentType, fileSize)
	if err != nil {
		return 0, fmt.Errorf("upload file to storage: %w", err)
	}

	obj := models.FileCU{
		Key: pointer.Of(fileKey),
	}

	id, err := uc.srv.File.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	return id, nil
}
