package repo

import (
	"bytes"
	"context"
	"fmt"

	m "back/infra/minio"

	"github.com/minio/minio-go/v7"
)

type RepoMinio interface {
	UploadImage(ctx context.Context, file []byte, path string, mimeType string) (string, error)
}

type Minio struct {
	minio *m.MinioClient
}

func NewMinio(minioClient *m.MinioClient) RepoMinio {
	return &Minio{
		minio: minioClient,
	}
}

func (repo *Minio) UploadImage(ctx context.Context, file []byte, path string, mimeType string) (string, error) {
	fmt.Println(path)
	_, err := repo.minio.Client.PutObject(
		ctx,
		repo.minio.BucketName,
		path,
		bytes.NewReader(file),
		int64(len(file)),
		minio.PutObjectOptions{ContentType: mimeType},
	)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s/%s/%s", repo.minio.Endpoint, repo.minio.BucketName, path)
	return url, nil
}
