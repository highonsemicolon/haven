package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/repo-handler-service/handlers"
)

func main() {
	r := gin.Default()

	repoHandler := handlers.NewRepoHandler()

	r.POST("/projects", repoHandler.CreateRepo)

	r.Run("localhost:8080")
}
