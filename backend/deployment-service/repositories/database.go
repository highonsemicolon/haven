package repositories

import (
	"database/sql"
	"log"

	"gorm.io/gorm"
)

func ConnectDatabase(dialector gorm.Dialector, config *gorm.Config, models ...interface{}) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		log.Fatalf("failed to connect database: %+v", err)
	}
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
	}

	log.Println("Database connected")
	return db, sqlDB
}
