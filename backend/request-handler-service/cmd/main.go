package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/request-handler-service/handlers"
	"github.com/onkarr19/haven/request-handler-service/services"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(-1, c.Errors.Last())
	}
}

func main() {
	r := gin.Default()
	r.Use(ErrorHandler)

	requestService := services.NewRequestService()
	requestHandler := handlers.NewRequestHandler(requestService)

	r.GET("/", requestHandler.GetDeployment)

	r.Run("localhost:8080")
}
