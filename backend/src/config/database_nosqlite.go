//go:build !sqlite

package config

import (
	"fmt"

	"gorm.io/gorm"
)

// connectSQLite is unavailable in the default (MySQL-only) build so that the
// CGO-based go-sqlite3 driver is never linked into production images.
func connectSQLite(cfg DBConfig) (*gorm.DB, error) {
	return nil, fmt.Errorf("binary built without sqlite support; rebuild with -tags sqlite")
}
