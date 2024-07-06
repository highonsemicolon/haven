package services

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/highonsemicolon/haven/deployment-service/models"
	"github.com/highonsemicolon/haven/deployment-service/repositories"
)

type DeploymentService interface {
	CreateDeployment(context.Context, *models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
	StreamLogs(context.Context, string) (<-chan string, error)
	AddConnection(string, *websocket.Conn)
	RemoveConnection(string)
	ListenForCloseMessages(string)
}

type deploymentService struct {
	deploymentRepo repositories.DeploymentRepository
	connections    map[string]*websocket.Conn
	mu             sync.Mutex
}

func NewDeploymentService(deploymentRepo repositories.DeploymentRepository) DeploymentService {
	return &deploymentService{deploymentRepo: deploymentRepo, connections: make(map[string]*websocket.Conn)}
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

func (s *deploymentService) AddConnection(id string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[id] = conn
}

func (s *deploymentService) RemoveConnection(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, ok := s.connections[id]; ok {
		conn.Close()
		delete(s.connections, id)
	}
}

func (s *deploymentService) ListenForCloseMessages(closingChannel string) {
	msgChan := s.deploymentRepo.Subscribe(closingChannel)
	for msg := range msgChan {
		s.RemoveConnection(msg)
	}
}
