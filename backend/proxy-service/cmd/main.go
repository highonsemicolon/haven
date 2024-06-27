package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/proxy-service/handlers"
	"github.com/onkarr19/haven/proxy-service/repositories"
	"github.com/onkarr19/haven/proxy-service/services"
	"github.com/sirupsen/logrus"
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

	base_path := os.Getenv("BASE_PATH")

	s3Repo := repositories.NewProxyRepository(base_path)
	requestService := services.NewProxyService(s3Repo, logrus.New())
	requestHandler := handlers.NewProxyHandler(requestService)

	r.NoRoute(requestHandler.HandleProxy)

	r.Run("localhost:8080")
}
