package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/server"
	serveraccrual "github.com/mi4r/gophermart/internal/server/accrual"
	"github.com/mi4r/gophermart/internal/storage"
	"github.com/mi4r/gophermart/lib/logger"
)

// Documentation: https://github.com/swaggo/swag
// @title Accrual System
// @version 1.0
// @description Swagger for Gopher Market API
// @host localhost:8080
// @BasePath /
func main() {
	config := config.NewAccSysConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorageAccrual(config.DriverType, config.StoragePath)
	core := server.NewServer(
		server.Config{
			ServiceName: server.GophermartName,
			Listen:      config.ListenAddr,
			MigrDirName: config.MigrDirName,
		},
	)
	service := serveraccrual.NewAccrualSystem(core)

	// Configure
	service.SetRoutes()
	service.SetStorage(storage)
	go service.Server.Start()

	// Канал для сигнала завершения
	done := make(chan struct{})

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	slog.Debug("received signal", slog.String("signal", sig.String()))
	close(done)
	service.Server.Shutdown()
}
