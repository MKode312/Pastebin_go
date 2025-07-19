package minio

import (
	"context"
	"text_sharing/internal/config"
	"text_sharing/internal/lib/utils"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client interface {
	InitMinio() error
	CreateOne(file utils.FileDataType, expires time.Duration) (string, string, error)
	GetOne(fileName string, date time.Time) (string, error)
}

type minioClient struct {
	mc *minio.Client
}

func NewMinioClient() Client {
	return &minioClient{}
}

func (m *minioClient) InitMinio() error {
	ctx := context.Background()

	client, err := minio.New(config.CfgMinio.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.CfgMinio.MinioRootUser, config.CfgMinio.MinioRootPassword, ""),
		Secure: config.CfgMinio.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	m.mc = client

	exists, err := m.mc.BucketExists(ctx, config.CfgMinio.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, config.CfgMinio.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
