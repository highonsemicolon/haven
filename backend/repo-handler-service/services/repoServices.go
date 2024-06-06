package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/repositories"
)

type RepoService interface {
	CreateRepo(repo *models.Repo) error
	IsUniqueRepo(name string) bool
}

type repoService struct {
	repoRepo repositories.RepoRepository
}

func NewRepoService(repoRepo repositories.RepoRepository) RepoService {
	return &repoService{repoRepo: repoRepo}
}

func (s *repoService) CreateRepo(repo *models.Repo) error {

	bucketName := "s3-haven--use1-az4--x-s3"
	objectKey := "abc.zip"

	presignedURL, err := getPresignedURL(bucketName, objectKey)
	if err != nil {
		return err
	}

	repo.PresignedURL = presignedURL

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, _ := cli.ContainerCreate(ctx, &container.Config{
		Image:      "nodewithgit",
		WorkingDir: "/app",
		Cmd: []string{"sh", "-c", fmt.Sprintf(`git clone %s . && npm install && npm run build && zip -r build-artifacts.zip build/* && 
			curl --upload-file build-artifacts.zip "%s"`, repo.GitURL, repo.PresignedURL)},
	}, nil, nil, nil, "")

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
	// if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
	// 	panic(err)
	// }

	//fmt.Println("Container removed successfully")
	s.repoRepo.CreateRepo(repo)

	return nil

}

func (s *repoService) IsUniqueRepo(name string) bool {
	// TODO: check if the repository name already does not exists
	return true
}
