package workeraccrual

import (
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
func NewWorker(id int, taskCh chan Task) *Worker {
	return &Worker{
		ID:     id,
		TaskCh: taskCh,
		QuitCh: make(chan struct{}),
	}
}

// Start запускает воркера
func (w *Worker) Start() {
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

func (w *Worker) SetStorage(storage storage.StorageAccrualSystem) {
	w.Storage = storage
	if err := w.Storage.Open(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (w *Worker) AddTask(task Task) {
	w.TaskCh <- task
}

func (w *Worker) Execute(task Task) error {
	slog.Debug("worker calculating accrual...", slog.String("order", task.Order.Order))

	if err := w.Storage.OrderRegUpdateStatus(storagedefault.StatusProcessing, task.Order.Order); err != nil {
		return err
	}

	rewards, err := w.Storage.RewardReadAll()
	if err != nil {
		return err
	}
	slog.Debug("rewards", slog.Any("rewards", rewards))

	var accrual float64
	for _, good := range task.Order.Goods {
		for _, reward := range rewards {
			var found bool
			if strings.Contains(good.Description, reward.Match) {
				slog.Debug("match one",
					slog.String("description", good.Description),
					slog.Float64("price", good.Price),
					slog.Float64("reward", reward.Reward),
					slog.String("type", string(reward.RewardType)),
				)
				accrual += calculateReward(good.Price, reward.Reward, reward.RewardType)
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

	if err := w.Storage.OrderRegUpdateOne(order); err != nil {
		return err
	}

	return nil
}

func calculateReward(price, reward float64, rewardType storageaccrual.RewardType) float64 {
	switch rewardType {
	case storageaccrual.RewardTypePercent:
		return (price / 100) * reward
	case storageaccrual.RewardTypePt:
		return reward
	default:
		return 0
	}
}
