package configs

import (
	"fmt"
)

type AwsConfig struct {
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	baseURL     string
}

func NewAwsConfig() *AwsConfig {
	cfg := &AwsConfig{
		S3Region:    getEnv("S3_REGION", "eu-central-1"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "key"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "secret"),
		S3Bucket:    getEnv("S3_BUCKET", "houses"),
	}
	cfg.baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", cfg.S3Bucket, cfg.S3Region)
	return cfg
}

func (c *AwsConfig) BaseURL() string {
	return c.baseURL
}

func (c *AwsConfig) AwsS3URL(key string) string {
	if key == "" {
		return ""
	}

	if len(key) > 4 && key[:4] == "http" {
		return key
	}

	if key[0] == '/' {
		key = key[1:]
	}

	return c.baseURL + key
}
