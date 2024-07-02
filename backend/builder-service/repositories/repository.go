package repositories

import (
	"context"

	"github.com/highonsemicolon/haven/builder-service/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	Pop(context.Context) (string, error)
	Push(context.Context, []byte) error
	Publish(context.Context, string, []byte) error

	GetDeploymentByName(name string) (*models.Builder, error)
	CreateDeployment(deployment *models.Builder) error
}

type brokerRepository struct {
	rds         *redis.Client
	inputQueue  string
	outputQueue string
	db          *gorm.DB
}

func NewBrokerRepository(db *gorm.DB, rds *redis.Client, inputQueue, outputQueue string) Repository {
	return &brokerRepository{
		db:          db,
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

func (r *brokerRepository) Push(ctx context.Context, message []byte) error {
	return r.rds.RPush(ctx, r.outputQueue, message).Err()
}

func (r *brokerRepository) Publish(ctx context.Context, channel string, message []byte) error {
	return r.rds.Publish(ctx, channel, message).Err()
}

func (r *brokerRepository) GetDeploymentByName(name string) (*models.Builder, error) {
	deployment := &models.Builder{}
	err := r.db.Where("name = ?", name).First(deployment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *brokerRepository) CreateDeployment(deployment *models.Builder) error {
	return r.db.Create(deployment).Error
}
