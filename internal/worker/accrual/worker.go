package workeraccrual

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/mi4r/gophermart/internal/storage"
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
)

// Думаю можно сделать пул воркеров
type Worker struct {
	ID      int           // ID воркера
	TaskCh  chan Task     // Канал для получения задач
	QuitCh  chan struct{} // Канал для завершения работы воркера
	Storage storage.StorageAccrualSystem
}

// NewWorker создает новый экземпляр воркера
func NewWorker(id int, taskCh chan Task, storage storage.StorageAccrualSystem) *Worker {
	return &Worker{
		ID:      id,
		TaskCh:  taskCh,
		QuitCh:  make(chan struct{}),
		Storage: storage,
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	if err := w.Storage.Open(); err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}
	go func() {
		for {
			select {
			case task := <-w.TaskCh:
				// Выполнение задачи
				if err := w.Execute(task); err != nil {
					slog.Error(err.Error(), slog.Int("id", w.ID))
				}
				slog.Debug("worker executed", slog.Int("id", w.ID))
			case <-w.QuitCh:
				// Завершение работы воркера
				slog.Debug("worker stopped", slog.Int("id", w.ID))
				return
			}
		}
	}()
}

// Stop останавливает воркера
func (w *Worker) Stop() {
	w.Storage.Close()
	go func() {
		w.QuitCh <- struct{}{}
	}()
}

func (w *Worker) AddTask(task Task) {
	w.TaskCh <- task
}

func (w *Worker) Execute(task Task) error {
	ctx := context.Background()
	slog.Debug("worker calculating accrual...", slog.String("order", task.Order.Order))

	if err := w.Storage.OrderRegUpdateStatus(ctx, storagedefault.StatusProcessing, task.Order.Order); err != nil {
		return err
	}

	rewards, err := w.Storage.RewardReadAll(ctx)
	if err != nil {
		return err
	}

	var accrual float64
	for _, good := range task.Order.Goods {
		for _, reward := range rewards {
			var found bool
			if strings.Contains(good.Description, reward.Match) {
				accrual += w.Calculate(good.Price, reward.Reward, reward.RewardType)
				found = true
			}
			// Если найдено 1 совпадение, то не продолжаем поиск
			if found {
				continue
			}
		}
	}

	order := storagedefault.Order{
		Number:  task.Order.Order,
		Status:  storagedefault.StatusProcessed,
		Accrual: accrual,
	}

	if err := w.Storage.OrderRegUpdateOne(ctx, order); err != nil {
		return err
	}

	return nil
}

func (w *Worker) Calculate(price, reward float64, rewardType storageaccrual.RewardType) float64 {
	switch rewardType {
	case storageaccrual.RewardTypePercent:
		return (price / 100) * reward
	case storageaccrual.RewardTypePt:
		return reward
	default:
		return 0
	}
}
