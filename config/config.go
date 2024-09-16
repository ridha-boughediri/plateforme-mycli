// config/config.go
package config

import "os"

type Config struct {
    AwsAccessKeyID     string
    AwsSecretAccessKey string
    AwsRegion          string
    S3Endpoint         string
}

func New() *Config {
    return &Config{
        AwsAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
        AwsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
        AwsRegion:          os.Getenv("AWS_REGION"),
        S3Endpoint:         os.Getenv("S3_ENDPOINT"),
    }
}
