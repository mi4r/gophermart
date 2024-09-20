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
)

func main() {
	config := config.NewAccSysConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorage(config.DriverType, config.StoragePath)
	core := server.NewServer(
		server.Config{
			ServiceName: server.GophermartName,
			Listen:      config.ListenAddr,
		}, storage,
	)
	service := server.NewAccrualSystem(core)
	service.SetRoutes()
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
