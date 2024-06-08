package services

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/repositories"
)

type RepoService interface {
	CreateRepo(repo *models.Repo) error
}

type repoService struct {
	repoRepo repositories.RepoRepository
}

func NewRepoService(repoRepo repositories.RepoRepository) RepoService {
	return &repoService{repoRepo: repoRepo}
}

func (s *repoService) CreateRepo(repo *models.Repo) error {

	// Check if repo.Name already exists
	if _, err := s.repoRepo.GetRepoByName(repo.Name); err == nil {
		return fmt.Errorf("project with name %s already exists", repo.Name)
	}

	// Generate a presigned URL
	presignedURL, err := putPresignURL("zip-builds/" + repo.Name)
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
            curl --upload-file build-artifacts.zip "%s"`, repo.GitURL, presignedURL),
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

	fmt.Println("Build completed successfully")

	// Remove the container
	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		panic(err)
	}

	s.repoRepo.CreateRepo(repo)

	return nil

}
