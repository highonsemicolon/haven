package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/request-handler-service/handlers"
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

	requestHandler := handlers.NewRequestHandler()

	r.GET("", requestHandler.GetDeployment)

	r.Run("localhost:8080")
}
