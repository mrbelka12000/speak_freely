package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
)

type (
	Storage struct {
		client *minio.Client
		bucket string
	}
)

func Connect(cfg config.Config) (*Storage, error) {

	minioClient, err := minio.New(cfg.MinIOAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %v", err)
	}

	return &Storage{
		client: minioClient,
		bucket: cfg.MinIOBucket,
	}, nil
}

func (s *Storage) UploadFile(ctx context.Context, file io.Reader, objectName string, fileSize int64) error {
	s.client.PutObject(ctx, s.bucket, objectName, file, fileSize, minio.PutObjectOptions{})
	return nil
}
