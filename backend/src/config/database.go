package config

import (
	"fmt"
	"log"
	"time"

	"groundTurn/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect opens the database and auto-migrates schema. The docker deployment
// uses MySQL; local/offline runs opt into SQLite via -tags sqlite and
// DB_DRIVER=sqlite.
func Connect(cfg DBConfig) (*gorm.DB, error) {
	if cfg.Driver == "sqlite" {
		return connectSQLite(cfg)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := migrate(db); err != nil {
		return nil, err
	}

	DB = db
	log.Printf("database connected driver=%s", cfg.Driver)
	return db, nil
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.FlightTurnaround{},
		&models.GroundTask{},
		&models.GroundResource{},
		&models.ResourceBooking{},
		&models.DelayEvent{},
		&models.AuditLog{},
	); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
