package workermart

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mi4r/gophermart/internal/storage"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
)

type Worker struct {
	ID             int          // ID воркера
	TickerCh       *time.Ticker // Канал для получения задач
	AccrualAddress string
	Storage        storage.StorageGophermart
}

// NewWorker создает новый экземпляр воркера
func NewWorker(id int, tickerCh *time.Ticker, accrualAddress string) *Worker {
	return &Worker{
		ID:             id,
		TickerCh:       tickerCh,
		AccrualAddress: accrualAddress,
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	slog.Debug("start timer")
	for range w.TickerCh.C {
		slog.Debug("start worker")
		w.Execute()
	}
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
	ctx := context.Background()
	orderNumbers, err := w.Storage.UserOrderReadAllNumbers(ctx)
	if err != nil {
		return err
	}

	slog.Debug("fetch orders", slog.Any("orders", orderNumbers))

	if len(orderNumbers) == 0 {
		return nil
	}

	var orders []storagedefault.Order

	for _, num := range orderNumbers {
		err = w.Storage.UserOrderUpdateStatus(ctx, num, storagedefault.StatusProcessing)
		if err != nil {
			return err
		}

		address := fmt.Sprintf("%s/api/orders/%s", w.AccrualAddress, num)
		slog.Debug("fetch data from accrual", slog.String("address", address))
		resp, err := http.Get(address)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			w.TickerCh.Reset(60 * time.Second)
			return nil
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		slog.Debug("response body from accrual", slog.String("respBody", string(respBody)))
		var order storagedefault.Order
		err = json.Unmarshal(respBody, &order)
		if err != nil {
			return err
		}
		orders = append(orders, order)
	}
	slog.Debug("orders", slog.Any("orders", orders))
	err = w.Storage.UserOrderUpdateAll(ctx, orders)
	if err != nil {
		return err
	}
	return nil
}
