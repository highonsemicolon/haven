package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type BrokerRepository struct {
	rds         *redis.Client
	inputQueue  string
	outputQueue string
}

func NewBrokerRepository(rds *redis.Client, inputQueue, outputQueue string) *BrokerRepository {
	return &BrokerRepository{
		rds:         rds,
		inputQueue:  inputQueue,
		outputQueue: outputQueue,
	}
}

func (r *BrokerRepository) Pop(ctx context.Context) (string, error) {
	result, err := r.rds.BLPop(ctx, 0, r.inputQueue).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}

func (r *BrokerRepository) Push(ctx context.Context, message string) error {
	return r.rds.RPush(ctx, r.outputQueue, message).Err()
}
