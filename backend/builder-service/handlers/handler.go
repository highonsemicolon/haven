package handlers

import (
	"context"
	"encoding/json"

	"github.com/onkarr19/haven/builder-service/models"
	"github.com/onkarr19/haven/builder-service/services"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger  *logrus.Logger
	service services.BrokerService
}

func NewHandler(logger *logrus.Logger, service services.BrokerService) *Handler {
	return &Handler{logger: logger, service: service}
}

func (h *Handler) Listen(ctx context.Context) error {
	for {
		input, err := h.service.Receive(ctx)
		if err != nil {
			h.logger.Errorf("Error listening for messages: %v", err)
			continue
		}

		deployment := models.Builder{}
		if err := json.Unmarshal([]byte(input), &deployment); err != nil {
			h.logger.Errorf("Error unmarshaling JSON: %v", err)
			continue
		}

		existingDeployment, _ := h.service.GetDeploymentByName(deployment.Name)
		if existingDeployment != nil {
			h.logger.Errorf("deployment with namde %s already exists", deployment.Name)
			continue
		}

		if err := h.service.CreateBuild(&deployment); err != nil {
			h.logger.Errorf("Error while building: %s", err)
			continue
		}

		output, err := json.Marshal(deployment)
		if err != nil {
			h.logger.Errorf("Error marshalling job to JSON: %v", err)
			continue
		}

		if err := h.service.Send(ctx, output); err != nil {
			h.logger.Errorf("Error sending message: %v", err)
		}
	}
}
