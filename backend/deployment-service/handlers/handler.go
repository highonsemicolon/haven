package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/highonsemicolon/haven/deployment-service/models"
	"github.com/highonsemicolon/haven/deployment-service/services"
	"github.com/sirupsen/logrus"
)

type DeploymentHandler struct {
	deploymentService services.DeploymentService
	logger            *logrus.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewDeploymentHandler(deploymentService services.DeploymentService, logger *logrus.Logger) *DeploymentHandler {
	return &DeploymentHandler{deploymentService: deploymentService, logger: logger}
}

func (h *DeploymentHandler) CreateDeployment(c *gin.Context) {
	var deployment models.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		h.logger.Errorf("failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	existingDeployment, _ := h.deploymentService.GetDeploymentByName(deployment.Name)
	if existingDeployment != nil {
		h.logger.Errorf("deployment with name %s already exists", deployment.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("deployment with name %s already exists", deployment.Name).Error()})
		return
	}

	deployment.ID = uuid.New()
	if deployment.Branch == "" {
		deployment.Branch = "main"
	}

	if err := h.deploymentService.CreateDeployment(c.Request.Context(), &deployment); err != nil {
		h.logger.Errorf("failed to create deployment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deployment"})
		return
	}

	c.JSON(http.StatusCreated, deployment)
}

func (h *DeploymentHandler) GetDeployment(c *gin.Context) {

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deployment name is required"})
		return
	}

	deployment, err := h.deploymentService.GetDeploymentByName(name)
	if err != nil {
		h.logger.Errorf("failed to get deployment by name %s: %v", name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get deployment"})
		return
	}

	if deployment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func (h *DeploymentHandler) HandleWebSocket(c *gin.Context) {
	deploymentID := c.Param("id")
	if deploymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deployment ID is required"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Errorf("failed to upgrade connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer ws.Close()

	ctx := c.Request.Context()
	logs, err := h.deploymentService.StreamLogs(ctx, deploymentID)
	if err != nil {
		h.logger.Errorf("failed to stream logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream logs"})
		return
	}

	for log := range logs {
		if err := ws.WriteMessage(websocket.BinaryMessage, []byte(log)); err != nil {
			h.logger.Errorf("failed to write message: %v", err)
			return
		}
	}
}
