package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/mrbelka12000/speak_freely/pkg/config"
)

type (
	Storage struct {
		client *minio.Client
		addr   string
		bucket string
	}
)

func Connect(cfg config.Config) (*Storage, error) {

	minioClient, err := minio.New(cfg.MinIOAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %v", err)
	}

	return &Storage{
		client: minioClient,
		bucket: cfg.MinIOBucket,
		addr:   cfg.MinIOAddr,
	}, nil
}

func (s *Storage) UploadFile(ctx context.Context, file io.Reader, objectName, contentType string, fileSize int64) (string, error) {
	info, err := s.client.PutObject(ctx, s.bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload file: %v", err)
	}

	return info.Key, nil
}
