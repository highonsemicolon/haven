package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/repo-handler-service/handlers"
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/repositories"
	"github.com/onkarr19/haven/repo-handler-service/services"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(-1, c.Errors.Last())
	}
}

func ConnectDatabase(sql gorm.Dialector, config *gorm.Config, models *models.Repo) *gorm.DB {
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
	r := gin.Default()
	r.Use(ErrorHandler)

	sql := sqlite.Open("test.db")
	db := ConnectDatabase(sql, &gorm.Config{}, &models.Repo{})

	repoRepository := repositories.NewRepoRepository(db)
	repoService := services.NewRepoService(repoRepository)
	repoHandler := handlers.NewRepoHandler(repoService)

	r.POST("/projects", repoHandler.CreateRepo)

	r.Run("localhost:8080")
}
