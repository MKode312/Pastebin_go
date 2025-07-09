package minio

import (
	"context"
	"text_sharing/internal/config"
	"text_sharing/internal/lib/utils"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio() error                                                                 // Метод для инициализации подключения к Minio
	CreateOne(file utils.FileDataType, expires time.Duration) (string, string, error) // Метод для создания одного объекта в бакете Minio
	GetOne(fileName string, date time.Time) (string, error)                           // Метод для получения одного объекта из бакета Minio                             // Метод для удаления одного объекта из бакета Minio
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc *minio.Client // Клиент Minio
}

// NewMinioClient создает новый экземпляр Minio Client
func NewMinioClient() Client {
	return &minioClient{} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

// InitMinio подключается к Minio и создает бакет, если не существует
// Бакет - это контейнер для хранения объектов в Minio. Он представляет собой пространство имен, в котором можно хранить и организовывать файлы и папки.
func (m *minioClient) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.CfgMinio.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.CfgMinio.MinioRootUser, config.CfgMinio.MinioRootPassword, ""),
		Secure: config.CfgMinio.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	m.mc = client

	// Проверка наличия бакета и его создание, если не существует
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
