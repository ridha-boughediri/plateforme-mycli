package config

import (
	"os"
)

func GetS3URL() string {
	return os.Getenv("S3_URL")
}

func GetAccessKey() string {
	return os.Getenv("S3_ACCESS_KEY")
}

func GetSecretKey() string {
	return os.Getenv("S3_SECRET_KEY")
}
