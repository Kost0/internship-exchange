package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client        *minio.Client
	publicBucket  string
	privateBucket string
}

func NewMinioStorage(addr, user, pass string) (*MinioStorage, error) {
	client, err := minio.New(addr, &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{
		client:        client,
		publicBucket:  "avatars",
		privateBucket: "resumes",
	}, nil
}

func (s *MinioStorage) EnsureBuckets(ctx context.Context) error {
	buckets := []string{s.publicBucket, "logos", s.privateBucket}

	for _, bucket := range buckets {
		exists, err := s.client.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}
		if !exists {
			if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return err
			}
		}
	}

	publicPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`

	for _, bucket := range []string{s.publicBucket, "logos"} {
		policy := fmt.Sprintf(publicPolicy, bucket)
		if err := s.client.SetBucketPolicy(ctx, bucket, policy); err != nil {
			return err
		}
	}

	return nil
}

func (s *MinioStorage) UploadAvatar(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("%s/avatar", userID)
	_, err := s.client.PutObject(ctx, s.publicBucket, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/%s/%s", s.publicBucket, objectName), nil
}

func (s *MinioStorage) UploadLogo(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("%s/logo", userID)
	_, err := s.client.PutObject(ctx, "logos", objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/logos/%s", objectName), nil
}

func (s *MinioStorage) UploadResume(ctx context.Context, userID string, data []byte) (string, error) {
	objectName := fmt.Sprintf("%s/resume.pdf", userID)
	_, err := s.client.PutObject(ctx, s.privateBucket, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/pdf"},
	)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (s *MinioStorage) GetResumePresignedURL(ctx context.Context, userID string) (string, error) {
	objectName := fmt.Sprintf("%s/resume.pdf", userID)
	url, err := s.client.PresignedGetObject(ctx, s.privateBucket, objectName, time.Hour, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
