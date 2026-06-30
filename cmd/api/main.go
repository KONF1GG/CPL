package main

import (
	"contralPlane/internal/config"
	"contralPlane/internal/logger"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development"
	}

	logger, err := logger.Init(mode)

	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config",
			zap.Error(err),
		)
	}

	logger.Info("config loaded successfully",
		zap.String("http_port", cfg.HTTPPort),
	)
}
