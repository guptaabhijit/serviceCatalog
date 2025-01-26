package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"serviceCatalog/config"
	"serviceCatalog/internal/models"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Get the underlying *sql.DB instance to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set up connection pooling
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)       // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)       // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime) // Maximum lifetime of a connection

	err = db.AutoMigrate(&models.Service{}, &models.Version{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
