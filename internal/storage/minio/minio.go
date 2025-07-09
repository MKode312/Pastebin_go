package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"text_sharing/internal/config"
	"text_sharing/internal/lib/utils"
	"text_sharing/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (m *minioClient) CreateOne(file utils.FileDataType, expires time.Duration) (string, string, error) {
	objectID := uuid.New().String()

	safeFileName := url.PathEscape(file.FileName)

	reader := bytes.NewReader(file.Data)

	_, err := m.mc.PutObject(context.Background(), config.CfgMinio.BucketName, safeFileName, reader, int64(len(file.Data)), minio.PutObjectOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error creating an object %s: %v", file.FileName, err)
	}

	url, err := m.mc.PresignedGetObject(context.Background(), config.CfgMinio.BucketName, safeFileName, expires, nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating an object %s: %v", file.FileName, err)
	}

	go func() {
		time.Sleep(expires)

		err := m.mc.RemoveObject(context.Background(), config.CfgMinio.BucketName, safeFileName, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("error deleting an object %s: %v\n", safeFileName, err)
		} else {
			fmt.Printf("object %s was successfully deleted\n", safeFileName)
		}
	}()

	return url.String(), objectID, nil
}

func (m *minioClient) GetOne(fileName string, date time.Time) (string, error) {
	ctx := context.Background()

	safeFileName := url.PathEscape(fileName)

	_, err := m.mc.StatObject(ctx, config.CfgMinio.BucketName, safeFileName, minio.StatObjectOptions{})
	if err != nil {
		return "", storage.ErrObjectExpired
	}

	object, err := m.mc.GetObject(ctx, config.CfgMinio.BucketName, safeFileName, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting an object: %w", err)
	}
	defer object.Close()

	filePath := "downloaded_file.txt"
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}

	_, err = io.Copy(file, object)
	if err != nil {
		file.Close()
		return "", fmt.Errorf("error copying data: %w", err)
	}
	file.Close()

	dataBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	err = os.Remove(filePath)

	content := string(dataBytes)

	return content, nil
}
