package storage

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"strings"
)

var logger = logging.GetLogger()

func Init(cfg *config.Config) *minio.Client {
	endpoint := "localhost:9000"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Successfully connected to MinIO")
	return minioClient
}

func uploadFile(file multipart.File, header *multipart.FileHeader, cfg *config.Config, minioClient *minio.Client) (string, error) {
	if file == nil || header == nil {
		return "", errors.New("invalid file")
	}
	defer file.Close()

	objectName := header.Filename
	contentType := header.Header.Get("Content-Type")
	bucketName := cfg.MinioBucket
	minioEndpoint := cfg.MinioEndpoint

	_, err := minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", errors.New("failed to uploadFile file")
	}

	fileURL := fmt.Sprintf("%s/%s/%s", minioEndpoint, bucketName, objectName)

	return fileURL, nil
}

func deleteFile(fileURL string, cfg *config.Config, minioClient *minio.Client) error {
	if fileURL == "" {
		return errors.New("missing file_url parameter")
	}

	prefix := cfg.MinioEndpoint + "/"
	if !strings.HasPrefix(fileURL, prefix) {
		return errors.New("invalid file_url format")
	}
	path := strings.TrimPrefix(fileURL, prefix)

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return errors.New("invalid file_url format")
	}

	bucketName := parts[0]
	objectName := parts[1]

	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.New("failed to deleteFile file")
	}

	return nil
}

type MinioStorage struct {
	client *minio.Client
	cfg    *config.Config
}

func (s *MinioStorage) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	return uploadFile(file, header, s.cfg, s.client)
}

func (s *MinioStorage) DeleteFileByURL(fileURL string) error {
	return deleteFile(fileURL, s.cfg, s.client)
}
