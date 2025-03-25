package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucketName string) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Явно отключаем SSL
	})
	if err != nil {
		return nil, err
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucketName,
		Endpoint:   endpoint,
	}, nil
}
