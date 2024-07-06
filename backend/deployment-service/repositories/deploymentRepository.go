package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/highonsemicolon/haven/deployment-service/models"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DeploymentRepository interface {
	CreateDeployment(context.Context, *models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
	StreamLogs(context.Context, string) (<-chan string, error)
	Subscribe(string) <-chan string
}

type deploymentRepository struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewDeploymentRepository(db *gorm.DB, rds *redis.Client) DeploymentRepository {
	return &deploymentRepository{db: db, rds: rds}
}

func (r *deploymentRepository) CreateDeployment(ctx context.Context, deployment *models.Deployment) error {
	job, err := json.Marshal(deployment)
	if err != nil {
		return errors.Wrap(err, "error marshalling job to JSON")
	}

	inputQueue := os.Getenv("BUILDER_QUEUE")

	if _, err := r.rds.RPush(ctx, inputQueue, job).Result(); err != nil {
		return errors.Wrap(err, "error pushing job to Redis")
	}
	return r.db.Create(deployment).Error
}

func (r *deploymentRepository) GetDeploymentByName(name string) (*models.Deployment, error) {
	deployment := &models.Deployment{}
	err := r.db.Where("name = ?", name).First(deployment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *deploymentRepository) StreamLogs(ctx context.Context, id string) (<-chan string, error) {
	logChannel := make(chan string)
	go func() {
		defer close(logChannel)
		pubsub := r.rds.Subscribe(ctx, id)
		defer pubsub.Close()

		for msg := range pubsub.Channel() {
			logChannel <- msg.Payload
		}
	}()
	return logChannel, nil
}

func (r *deploymentRepository) Subscribe(channel string) <-chan string {
	ctx := context.Background()
	pubsub := r.rds.Subscribe(ctx, channel)
	ch := make(chan string)

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("error receiving message:", err)
				continue
			}

			ch <- msg.Payload
		}
	}()

	return ch
}
