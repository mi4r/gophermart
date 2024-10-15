package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mi4r/gophermart/docs/gophermart"
	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/server"
	servermart "github.com/mi4r/gophermart/internal/server/gophermart"
	"github.com/mi4r/gophermart/internal/storage"
	workermart "github.com/mi4r/gophermart/internal/worker/gophermart"
	"github.com/mi4r/gophermart/lib/logger"
)

// Documentation: https://github.com/swaggo/swag
// @title Gophermart
// @version 1.0
// @description Swagger for Gopher Market API
// @host localhost:8080
// @BasePath /
func main() {
	config := config.NewGophermartConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorageGophermart(config.DriverType, config.StoragePath)
	core := server.NewServer(
		server.Config{
			ServiceName: server.GophermartName,
			Listen:      config.ListenAddr,
			SecretKey:   config.SecretKey,
		},
	)

	tickerCh := time.NewTicker(config.TickerTime)
	worker := workermart.NewWorker(1, tickerCh, config.AccrualSystemAddress)
	worker.SetStorage(storage)
	service := servermart.NewGophermart(core)
	// Configure
	service.SetRoutes()
	service.SetStorage(storage)
	go service.Server.Start()
	go worker.Start()

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	slog.Debug("received signal", slog.String("signal", sig.String()))
	worker.Stop()
	service.Server.Shutdown()
}
