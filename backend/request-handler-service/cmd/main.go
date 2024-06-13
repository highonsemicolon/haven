package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/request-handler-service/handlers"
	"github.com/onkarr19/haven/request-handler-service/repositories"
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

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("DEPLOYMENT_BUCKET")

	s3Repo, _ := repositories.NewS3Repository(region, bucket)
	requestService := services.NewRequestService(s3Repo)
	requestHandler := handlers.NewRequestHandler(requestService)

	r.NoRoute(requestHandler.GetDeployment)

	r.Run("localhost:8080")
}
