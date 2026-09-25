package config

import (
	"fmt"
	"os"
)

// DBConfig is assembled from environment variables (docker-compose/.env),
// keeping config intentionally spread across .env, compose and this package.
type DBConfig struct {
	Driver string // mysql | sqlite
	DSN    string
}

func GetDBConfig() DBConfig {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}
	if driver == "sqlite" {
		path := os.Getenv("DB_SQLITE_PATH")
		if path == "" {
			path = "ground-turn.db"
		}
		return DBConfig{Driver: "sqlite", DSN: path}
	}

	user := envOr("DB_USER", "ground")
	pass := envOr("DB_PASSWORD", "ground")
	host := envOr("DB_HOST", "db")
	port := envOr("DB_PORT", "3306")
	name := envOr("DB_NAME", "ground_turn")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)
	return DBConfig{Driver: "mysql", DSN: dsn}
}

// ServerPort returns the backend listen port (container internal is 3000).
func ServerPort() string {
	return envOr("SERVER_PORT", "3000")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
