package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/services"
	"github.com/sirupsen/logrus"
)

type DeploymentHandler struct {
	deploymentService services.DeploymentService
	logger            *logrus.Logger
}

func NewDeploymentHandler(deploymentService services.DeploymentService, logger *logrus.Logger) *DeploymentHandler {
	return &DeploymentHandler{deploymentService: deploymentService, logger: logger}
}

func (h *DeploymentHandler) CreateDeployment(c *gin.Context) {
	var deployment models.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	deployment.ID = uuid.New()

	if err := h.deploymentService.CreateDeployment(&deployment); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
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
