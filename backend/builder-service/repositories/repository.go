package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Pop(ctx context.Context) (string, error)
	Push(ctx context.Context, message string) error
}

type brokerRepository struct {
	rds         *redis.Client
	inputQueue  string
	outputQueue string
}

func NewBrokerRepository(rds *redis.Client, inputQueue, outputQueue string) Repository {
	return &brokerRepository{
		rds:         rds,
		inputQueue:  inputQueue,
		outputQueue: outputQueue,
	}
}

func (r *brokerRepository) Pop(ctx context.Context) (string, error) {
	result, err := r.rds.BLPop(ctx, 0, r.inputQueue).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}

func (r *brokerRepository) Push(ctx context.Context, message string) error {
	return r.rds.RPush(ctx, r.outputQueue, message).Err()
}
