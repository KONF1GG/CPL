package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort         string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
	ShutdownTimeout  time.Duration

	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	dbHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := requiredEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbUser, err := requiredEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPassword, err := requiredEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbName, err := requiredEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbMaxOpenConns, err := positiveIntEnv("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return nil, err
	}
	dbMaxIdleConns, err := positiveIntEnv("DB_MAX_IDLE_CONNS", 25)
	if err != nil {
		return nil, err
	}
	if dbMaxIdleConns > dbMaxOpenConns {
		return nil, fmt.Errorf("DB_MAX_IDLE_CONNS must be less than or equal to DB_MAX_OPEN_CONNS")
	}

	httpReadTimeout, err := durationEnv("HTTP_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, err
	}
	httpWriteTimeout, err := durationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}
	httpIdleTimeout, err := durationEnv("HTTP_IDLE_TIMEOUT", 120*time.Second)
	if err != nil {
		return nil, err
	}
	shutdownTimeout, err := durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}
	dbConnMaxLifetime, err := durationEnv("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}
	dbConnMaxIdleTime, err := durationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPPort:         envOrDefault("HTTP_PORT", "8080"),
		HTTPReadTimeout:  httpReadTimeout,
		HTTPWriteTimeout: httpWriteTimeout,
		HTTPIdleTimeout:  httpIdleTimeout,
		ShutdownTimeout:  shutdownTimeout,

		DBHost:            dbHost,
		DBPort:            dbPort,
		DBUser:            dbUser,
		DBPassword:        dbPassword,
		DBName:            dbName,
		DBSSLMode:         envOrDefault("DB_SSL_MODE", "disable"),
		DBMaxOpenConns:    dbMaxOpenConns,
		DBMaxIdleConns:    dbMaxIdleConns,
		DBConnMaxLifetime: dbConnMaxLifetime,
		DBConnMaxIdleTime: dbConnMaxIdleTime,
	}, nil
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not set", key)
	}
	return value, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func positiveIntEnv(key string, fallback int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return value, nil
}

func durationEnv(key string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return value, nil
}
