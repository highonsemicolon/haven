package services

import (
	"context"

	"github.com/onkarr19/haven/deployment-service/models"
	"github.com/onkarr19/haven/deployment-service/repositories"
)

type DeploymentService interface {
	CreateDeployment(context.Context, *models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
	StreamLogs(context.Context, string) (<-chan string, error)
}

type deploymentService struct {
	deploymentRepo repositories.DeploymentRepository
}

func NewDeploymentService(deploymentRepo repositories.DeploymentRepository) DeploymentService {
	return &deploymentService{deploymentRepo: deploymentRepo}
}

func (s *deploymentService) CreateDeployment(ctx context.Context, deployment *models.Deployment) error {
	return s.deploymentRepo.CreateDeployment(ctx, deployment)
}

func (s *deploymentService) GetDeploymentByName(name string) (*models.Deployment, error) {
	return s.deploymentRepo.GetDeploymentByName(name)
}

func (s *deploymentService) StreamLogs(ctx context.Context, id string) (<-chan string, error) {
	return s.deploymentRepo.StreamLogs(ctx, id)
}
