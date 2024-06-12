package repositories

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository interface {
	GetObject(context.Context, string) (io.ReadCloser, string, error)
}

type s3Repository struct {
	s3Client *s3.Client
	bucket   string
}

func NewS3Repository(region, bucket string) (S3Repository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return &s3Repository{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
	}, nil
}

func (p *s3Repository) GetObject(ctx context.Context, key string) (io.ReadCloser, string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}
	result, err := p.s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object from S3: %w", err)
	}

	return result.Body, *result.ContentType, nil
}
