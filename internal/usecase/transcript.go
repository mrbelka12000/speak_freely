package usecase

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

func (uc *UseCase) TranscriptCreate(
	ctx context.Context,
	b *bytes.Buffer,
	objectName,
	contentType string,
	fileSize int64,
	themeID int64,
	languageID int64,
) error {
	objectName = fmt.Sprintf("%d-%s", time.Now().UnixMilli(), objectName)

	fileKey, err := uc.storage.UploadFile(ctx, b, objectName, contentType, fileSize)
	if err != nil {
		return fmt.Errorf("upload file to storage: %w", err)
	}

	_ = fileKey
	return nil
}
