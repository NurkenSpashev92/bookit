package configs

import (
	"fmt"
)

type AwsConfig struct {
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func NewAwsConfig() *AwsConfig {
	return &AwsConfig{
		S3Region:    getEnv("S3_REGION", "eu-central-1"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "key"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "secret"),
		S3Bucket:    getEnv("S3_BUCKET", "houses"),
	}
}

func (c *AwsConfig) AwsS3URL(key string) string {
	if key == "" {
		return ""
	}

	if len(key) > 4 && key[:4] == "http" {
		return key
	}

	base := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", c.S3Bucket, c.S3Region)

	if key[0] == '/' {
		key = key[1:]
	}

	return base + "/" + key
}
