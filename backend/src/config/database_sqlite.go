//go:build sqlite

package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connectSQLite is compiled only with `-tags sqlite` because the go-sqlite3
// driver requires CGO. Production images stay CGO-free via the MySQL path.
func connectSQLite(cfg DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	DB = db
	log.Printf("database connected driver=sqlite dsn=%s", cfg.DSN)
	return db, nil
}
