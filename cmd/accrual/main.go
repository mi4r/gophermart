package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mi4r/gophermart/docs/accrual"
	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/server"
	serveraccrual "github.com/mi4r/gophermart/internal/server/accrual"
	"github.com/mi4r/gophermart/internal/storage"
	workeraccrual "github.com/mi4r/gophermart/internal/worker/accrual"
	"github.com/mi4r/gophermart/lib/logger"
)

// Documentation: https://github.com/swaggo/swag
// @title Accrual System
// @version 1.0
// @description Swagger for Gopher Market API
// @host localhost:8081
// @BasePath /
func main() {
	config := config.NewAccrualConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorageAccrual(config.DriverType, config.StoragePath)
	core := server.NewServer(
		server.Config{
			ServiceName: server.AccrualName,
			Listen:      config.ListenAddr,
			RateLimit:   config.RateLimit,
		},
	)
	// Канал для передачи задач
	taskCh := make(chan workeraccrual.Task)
	worker := workeraccrual.NewWorker(1, taskCh)

	service := serveraccrual.NewAccrualSystem(core, taskCh)

	// Configure
	service.SetRoutes()
	service.SetStorage(storage)
	worker.SetStorage(storage)
	go service.Server.Start()
	worker.Start()
	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	slog.Debug("received signal", slog.String("signal", sig.String()))

	// Закрываем канал задач
	close(taskCh)
	service.Server.Shutdown()
}
