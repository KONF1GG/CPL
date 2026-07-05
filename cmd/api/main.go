package main

import (
	"CPL/internal/config"
	"CPL/internal/database"
	"CPL/internal/handler"
	"CPL/internal/logger"
	"CPL/internal/middleware"
	"CPL/internal/repository"
	"CPL/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	vmHandler := handler.NewVMHandler(vmService)
	taskHandler := handler.NewTaskHandler(taskService)

	router := handler.NewRouter(vmHandler, taskHandler)
	httpHandler := middleware.Logging(logger, middleware.Recovery(logger, router))

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      httpHandler,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	go func() {
		logger.Info("http server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Info("server stopped")
}
