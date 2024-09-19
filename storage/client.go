package storage

import (
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

// InitMinioClient initialise le client MinIO
func InitMinioClient() error {
	endpoint := os.Getenv("Endpoint")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	MinioClient = minioClient
	return nil
}
