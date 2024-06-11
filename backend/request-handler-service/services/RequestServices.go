package services

import (
	"strings"
)

type RequestService interface {
	GetSubdomain(host string) string
}

type requestService struct {
}

func NewRequestService() RequestService {
	return &requestService{}
}

func (s *requestService) GetSubdomain(host string) string {

	parts := strings.Split(host, ".")

	if len(parts) > 2 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}
