package services

import (
	"context"
	"encoding/json"
	"os"

	"github.com/onkarr19/haven/deployment-service/models"
	"github.com/onkarr19/haven/deployment-service/repositories"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type DeploymentService interface {
	CreateDeployment(*models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
}

type deploymentService struct {
	deploymentRepo repositories.DeploymentRepository
	rds            *redis.Client
}

func NewDeploymentService(deploymentRepo repositories.DeploymentRepository, redis *redis.Client) DeploymentService {
	return &deploymentService{deploymentRepo: deploymentRepo, rds: redis}
}

func (s *deploymentService) CreateDeployment(deployment *models.Deployment) error {
	job, err := json.Marshal(deployment)
	if err != nil {
		return errors.Wrap(err, "error marshalling job to JSON")
	}

	inputQueue := os.Getenv("INOUT_QUEUE")

	if _, err := s.rds.RPush(context.Background(), inputQueue, job).Result(); err != nil {
		return errors.Wrap(err, "error pushing job to Redis")
	}

	if err := s.deploymentRepo.CreateDeployment(deployment); err != nil {
		return errors.Wrap(err, "error saving deployment to repository")
	}
	return nil
}

func (s *deploymentService) GetDeploymentByName(name string) (*models.Deployment, error) {
	deployment, err := s.deploymentRepo.GetDeploymentByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "error getting deployment by name")
	}
	return deployment, nil
}
