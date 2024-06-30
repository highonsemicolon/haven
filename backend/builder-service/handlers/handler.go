package handlers

import (
	"context"

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
		message, err := h.service.Receive(ctx)
		if err != nil {
			h.logger.Errorf("Error listening for messages: %v", err)
			continue
		}

		response := h.service.Process(ctx, message)
		h.logger.Printf("Processed message: %s", response)

		if err := h.service.Send(ctx, response); err != nil {
			h.logger.Printf("Error sending message: %v", err)
		}
	}
}
