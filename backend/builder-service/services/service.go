package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/onkarr19/haven/builder-service/models"
	"github.com/onkarr19/haven/builder-service/repositories"
	"github.com/pkg/errors"
)

type BrokerService interface {
	Receive(ctx context.Context) (string, error)
	Send(ctx context.Context, message []byte) error
	Process(ctx context.Context, message string) string

	GetDeploymentByName(name string) (*models.Builder, error)
	CreateBuild(deployment *models.Builder) error
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

func (s *brokerService) Send(ctx context.Context, message []byte) error {
	return s.brokerRepository.Push(ctx, message)
}

func (s *brokerService) Process(ctx context.Context, message string) string {
	return message
}

func (s *brokerService) CreateBuild(deployment *models.Builder) error {

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
	if err := s.brokerRepository.CreateDeployment(deployment); err != nil {
		return errors.Wrap(err, "error saving deployment to repository")
	}
	return nil
}

func (s *brokerService) GetDeploymentByName(name string) (*models.Builder, error) {
	deployment, err := s.brokerRepository.GetDeploymentByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "error getting deployment by name")
	}
	return deployment, nil
}
