package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/server"
	"github.com/mi4r/gophermart/internal/storage"
	"github.com/mi4r/gophermart/lib/logger"
)

func main() {
	config := config.NewServerConfig()
	logger.InitLogger(config.LogLevel)
	storage := storage.NewStorage(config.DriverType)
	server := server.NewServer(
		config, &storage,
	)

	wg := &sync.WaitGroup{}
	// Идея в том чтобы при получении сигнала счетчик становился 0
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.Start()
	}()

	go func() {
		// Канал для получения сигналов
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		// Ожидаем сигнал
		s := <-sigChan
		slog.Debug("A shutdown signal has been received", slog.String("signal", s.String()))
		wg.Done()
	}()

	wg.Wait()
}
