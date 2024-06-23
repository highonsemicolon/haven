package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func putPresignURL(objectName string) (string, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("unable to load AWS config: %w", err)
	}

	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("S3_BUCKET_NAME environment variable is not set")
	}

	// Create an S3 client
	s3client := s3.NewFromConfig(cfg)

	// Create a presign client
	presignClient := s3.NewPresignClient(s3client)

	// Generate a presigned URL
	presignedUrl, err := presignClient.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		},
		s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("unable to generate presigned URL: %w", err)
	}

	return presignedUrl.URL, nil
}
