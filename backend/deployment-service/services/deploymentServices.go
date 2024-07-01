package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/repositories"
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

	if _, err := s.rds.RPush(context.Background(), "builder-input", job).Result(); err != nil {
		return errors.Wrap(err, "error pushing job to Redis")
	}

	if err := s.deploymentRepo.CreateDeployment(deployment); err != nil {
		return errors.Wrap(err, "error saving deployment to repository")
	}
	return nil
}

func (s *deploymentService) CreateDeploymentbackup(deployment *models.Deployment) error {

	// Generate a presigned URL
	presignedURL, err := putPresignURL(deployment.Name)
	if err != nil {
		return errors.Wrap(err, "error generating presigned URL")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "error initializing Docker client")
	}
	defer cli.Close()

	ctx := context.Background()

	config := &container.Config{
		Image:      "nodewithgit",
		WorkingDir: "/app",
		Cmd: []string{
			"sh", "-c",
			fmt.Sprintf(`git clone -b "%s" "%s" . && npm install && npm run build && zip -r build-artifacts.zip build/* &&
		            curl --upload-file build-artifacts.zip "%s"`, deployment.Branch, deployment.GitURL, presignedURL),
		},
	}

	// Create a container
	resp, err := cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return errors.Wrap(err, "error creating Docker container")
	}
	defer func() {
		if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			errors.Wrap(err, "error removing Docker container")
		}
	}()

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return errors.Wrap(err, "error starting Docker container")
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "error waiting for Docker container to finish")
		}
	case <-statusCh:
	}

	// set the hosted URL for deployment
	deployment.HostedURL = "https://" + deployment.Name + ".haven.app"

	// Save the deployment to the repository
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
