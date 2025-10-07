package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

// DBConfig holds database configuration
type DBConfig struct {
	Addr            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoadDBConfig loads database configuration from environment variables
func LoadDBConfig() *DBConfig {
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	return &DBConfig{
		Addr:            getEnv("DB_ADDR", ""),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (config *DBConfig) GetConnectionString() string {
	if config.Addr == "" {
		panic("DB_ADDR environment variable is required")
	}
	return config.Addr
}

func NewDBConn(lc fx.Lifecycle) (*sql.DB, error) {
	config := LoadDBConfig()
	var db *sql.DB

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			db, err = sql.Open("pgx", config.GetConnectionString())
			if err != nil {
				return fmt.Errorf("failed to open database connection: %w", err)
			}

			db.SetMaxOpenConns(config.MaxOpenConns)
			db.SetMaxIdleConns(config.MaxIdleConns)
			db.SetConnMaxLifetime(config.ConnMaxLifetime)

			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database ping failed: %w", err)
			}

			log.Println("Database connection established successfully")
			return nil
		},

		OnStop: func(ctx context.Context) error {
			if db != nil {
				return db.Close()
			}
			return nil
		},
	})

	return db, nil
}
