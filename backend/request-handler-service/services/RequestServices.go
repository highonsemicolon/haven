package services

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/onkarr19/haven/request-handler-service/repositories"
)

type RequestService interface {
	GetSubdomain(host string) string
	GetDeploymentContent(ctx context.Context, subdomain string) (io.ReadCloser, string, error)
}

type requestService struct {
	s3Repo repositories.S3Repository
}

func NewRequestService(s3repo repositories.S3Repository) RequestService {
	return &requestService{s3Repo: s3repo}
}

func (s *requestService) GetSubdomain(host string) string {
	parts := strings.Split(host, ".")

	if len(parts) > 2 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

func (s *requestService) GetDeploymentContent(ctx context.Context, subdomain string) (io.ReadCloser, string, error) {
	key := subdomain + "/index.html"
	fmt.Println("Key: ", key)
	return s.s3Repo.GetObject(ctx, key)
}
