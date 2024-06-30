package services

import (
	"context"

	"github.com/onkarr19/haven/builder-service/repositories"
)

type BrokerService interface {
	Receive(ctx context.Context) (string, error)
	Send(ctx context.Context, message string) error
	Process(ctx context.Context, message string) string
}

type brokerService struct {
	brokerRepository repositories.Repository
}

func NewBrokerService(repo repositories.Repository) BrokerService {
	return &brokerService{brokerRepository: repo}
}

func (s *brokerService) Receive(ctx context.Context) (string, error) {
	return s.brokerRepository.Pop(ctx)
}

func (s *brokerService) Send(ctx context.Context, message string) error {
	return s.brokerRepository.Push(ctx, message)
}

func (s *brokerService) Process(ctx context.Context, message string) string {
	return message
}
