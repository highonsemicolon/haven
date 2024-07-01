package repositories

import (
	"database/sql"

	"gorm.io/gorm"
)

func ConnectDatabase(dialector gorm.Dialector, config *gorm.Config, models ...interface{}) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, nil, err
	}
	if err := db.AutoMigrate(models...); err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}
