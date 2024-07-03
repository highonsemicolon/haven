package handlers

import (
	"context"

	"github.com/highonsemicolon/haven/builder-service/services"
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

		var output []byte
		if output, err = h.service.PrepareBuild(input); err != nil {
			if err := h.service.Send(ctx, []byte(err.Error())); err != nil {
				h.logger.Errorf("Error sending message: %v", err)
			}
			continue
		}

		if err := h.service.Send(ctx, output); err != nil {
			h.logger.Errorf("Error sending message: %v", err)
		}
	}
}
