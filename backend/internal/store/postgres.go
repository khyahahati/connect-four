package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgres establishes a connection to Postgres using environment configuration.
func NewPostgres() (*sql.DB, error) {
	host, err := requiredEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	portStr, err := requiredEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	user, err := requiredEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	password, err := requiredEnv("DB_PASS")
	if err != nil {
		return nil, err
	}
	name, err := requiredEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	if _, err := strconv.Atoi(portStr); err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, portStr, user, password, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func requiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}
