package services

import (
	"context"

	"github.com/onkarr19/haven/builder-service/repositories"
)

type BrokerService struct {
	brokerRepository *repositories.BrokerRepository
}

func NewBrokerService(repo *repositories.BrokerRepository) *BrokerService {
	return &BrokerService{brokerRepository: repo}
}

func (s *BrokerService) Receive(ctx context.Context) (string, error) {
	return s.brokerRepository.Pop(ctx)
}

func (s *BrokerService) Send(ctx context.Context, message string) error {
	return s.brokerRepository.Push(ctx, message)
}

func (s *BrokerService) Process(ctx context.Context, message string) string {
	return message
}
