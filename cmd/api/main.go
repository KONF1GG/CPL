package main

import (
	"CPL/internal/config"
	"CPL/internal/database"
	"CPL/internal/logger"
	"CPL/internal/repository"
	"CPL/internal/service"
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

	db, err := database.NewDB(cfg)
	if err != nil {
		logger.Fatal("failed to connect to database",
			zap.Error(err),
		)
	}

	vmRepo := repository.NewVMRepository(db)
	taskRepo := repository.NewTaskRepo(db)
	txManager := repository.NewTransactionManager(db)

	vmService := service.NewVMService(vmRepo, taskRepo, txManager)
	taskService := service.NewTaskService(taskRepo)
	_ = vmService
	_ = taskService

	logger.Info("application started",
		zap.String("http_port", cfg.HTTPPort),
	)
}
