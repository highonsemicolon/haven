package services

import (
	"github.com/highonsemicolon/haven/proxy-service/repositories"
)

type ProxyService interface {
	GetObjectURL(string) string
}

type proxyService struct {
	proxyRepo repositories.ProxyRepository
}

func NewProxyService(s3repo repositories.ProxyRepository) ProxyService {
	return &proxyService{proxyRepo: s3repo}
}

func (s *proxyService) GetObjectURL(key string) string {
	return s.proxyRepo.GetObjectURL(key)
}
