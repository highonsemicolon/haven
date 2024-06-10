package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func (h *RequestHandler) GetDeployment(c *gin.Context) {

	fmt.Println("Inside GetDeployment")
	fmt.Printf("host: %+v\n", c.Request.Host)

	c.JSON(http.StatusOK, gin.H{})
}
