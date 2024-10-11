package workermart

import (
	"log/slog"
	"os"

	"github.com/mi4r/gophermart/internal/storage"
)

type Worker struct {
	ID      int           // ID воркера
	TaskCh  chan Task     // Канал для получения задач
	QuitCh  chan struct{} // Канал для завершения работы воркера
	Storage storage.StorageGophermart
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

func (w *Worker) SetStorage(storage storage.StorageGophermart) {
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
	return nil
}
