package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/server"
	"github.com/mi4r/gophermart/internal/storage"
	"github.com/mi4r/gophermart/lib/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mi4r/gophermart/docs"
)

// Documentation: https://github.com/swaggo/swag
// @title Gophermart
// @version 1.0
// @description Swagger for Gopher Market API
// @host localhost:8080
// @BasePath /
func main() {
	config := config.NewServerConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorage(config.DriverType, config.StoragePath)
	server := server.NewServer(
		config, &storage,
	)
	go server.Start()

	// Канал для сигнала завершения
	done := make(chan struct{})

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	slog.Debug("received signal", slog.String("signal", sig.String()))
	close(done)
	server.Shutdown()

}
