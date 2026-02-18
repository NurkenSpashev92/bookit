package aws

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsS3Client struct {
	client *s3.Client
	bucket string
	region string
}

func NewAwsS3Client(region, accessKey, secretKey, bucket string) (*AwsS3Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = false
		o.UseAccelerate = false
		o.UseARNRegion = true
	})

	return &AwsS3Client{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

func (s *AwsS3Client) Upload(ctx context.Context, bucket string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s/%d_%s",
		bucket,
		time.Now().UnixNano(),
		filepath.Base(file.Filename),
	)

	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *AwsS3Client) UploadCompressed(ctx context.Context, key string, body []byte, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *AwsS3Client) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucket, s.region)
	if len(key) > len(baseURL) && key[:len(baseURL)] == baseURL {
		key = key[len(baseURL):]
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}
