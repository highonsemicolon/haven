package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/url-management-service/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)
	err := r.Run("localhost:8083")
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
}
