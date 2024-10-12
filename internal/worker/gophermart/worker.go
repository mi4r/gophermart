package workermart

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/mi4r/gophermart/internal/storage"
)

type Worker struct {
	ID       int          // ID воркера
	TickerCh *time.Ticker // Канал для получения задач
	Storage  storage.StorageGophermart
}

// NewWorker создает новый экземпляр воркера
func NewWorker(id int, tickerCh *time.Ticker) *Worker {
	return &Worker{
		ID:       id,
		TickerCh: tickerCh,
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	go func() {
		for range w.TickerCh.C {
			w.Execute()
		}
	}()
}

// Stop останавливает воркера
func (w *Worker) Stop() {
	w.Storage.Close()
	w.TickerCh.Stop()
}

func (w *Worker) SetStorage(storage storage.StorageGophermart) {
	w.Storage = storage
	ctx := context.Background()
	if err := w.Storage.Open(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (w *Worker) Execute() error {
	return nil
}
