package services

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/repositories"
)

type DeploymentService interface {
	CreateDeployment(*models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
}

type deploymentService struct {
	deploymentRepo repositories.DeploymentRepository
}

func NewDeploymentService(deploymentRepo repositories.DeploymentRepository) DeploymentService {
	return &deploymentService{deploymentRepo: deploymentRepo}
}

func (s *deploymentService) CreateDeployment(deployment *models.Deployment) error {

	// Check if deployment.Name already exists
	if _, err := s.deploymentRepo.GetDeploymentByName(deployment.Name); err == nil {
		return fmt.Errorf("project with name %s already exists", deployment.Name)
	}

	// Generate a presigned URL
	presignedURL, err := putPresignURL(deployment.Name)
	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()

	config := &container.Config{
		Image:      "nodewithgit",
		WorkingDir: "/app",
		Cmd: []string{
			"sh", "-c",
			fmt.Sprintf(`git clone "%s" . && npm install && npm run build && zip -r build-artifacts.zip build/* &&
		            curl --upload-file build-artifacts.zip "%s"`, deployment.GitURL, presignedURL),
		},
	}

	// Create a container
	resp, err := cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return err
	}
	defer func() {
		if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			log.Println("Failed to remove container:", err)
		}
	}()

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	s.deploymentRepo.CreateDeployment(deployment)

	return nil

}

func (s *deploymentService) GetDeploymentByName(name string) (*models.Deployment, error) {
	return s.deploymentRepo.GetDeploymentByName(name)
}
