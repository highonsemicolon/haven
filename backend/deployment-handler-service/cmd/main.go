package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onkarr19/haven/deployment-handler-service/handlers"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/repositories"
	"github.com/onkarr19/haven/deployment-handler-service/services"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(-1, c.Errors.Last())
	}
}

func ConnectDatabase(sql gorm.Dialector, config *gorm.Config, models *models.Deployment) *gorm.DB {
	db, err := gorm.Open(sql, config)
	if err != nil {
		log.Fatalf("failed to connect database: %+v", err)
	}
	if err := db.AutoMigrate(models); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connected")
	return db
}

func main() {
	godotenv.Load()
	r := gin.Default()
	r.Use(ErrorHandler)

	sql := sqlite.Open("test.db")
	db := ConnectDatabase(sql, &gorm.Config{}, &models.Deployment{})

	deploymentRepository := repositories.NewDeploymentRepository(db)
	deploymentService := services.NewDeploymentService(deploymentRepository)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService)

	r.POST("/project", deploymentHandler.CreateDeployment)
	r.GET("/project/:name", deploymentHandler.GetDeployment)

	r.Run("localhost:8080")
}
