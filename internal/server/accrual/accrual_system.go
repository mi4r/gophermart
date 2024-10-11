package serveraccrual

import (
	"log/slog"

	"github.com/mi4r/gophermart/internal/server"
	"github.com/mi4r/gophermart/internal/storage"
	workeraccrual "github.com/mi4r/gophermart/internal/worker/accrual"
	"golang.org/x/time/rate"
)

type AccrualSystem struct {
	*server.Server
	taskCh      chan workeraccrual.Task
	storage     storage.StorageAccrualSystem
	rateLimiter *rate.Limiter
}

func NewAccrualSystem(server *server.Server, taskCh chan workeraccrual.Task) *AccrualSystem {
	return &AccrualSystem{
		taskCh: taskCh,
		Server: server,
		// 5 requests in 1 minute
		rateLimiter: rate.NewLimiter(rate.Limit(server.Config.RateLimit), 60),
	}
}

func (s *AccrualSystem) SetRoutes() {
	gAPI := s.Router.Group("/api")
	gAPI.GET("/orders/:number", s.ordersGetHandler, server.RateLimiterMiddleware(s.rateLimiter))
	gAPI.POST("/orders", s.ordersPostHandler)
	gAPI.POST("/goods", s.rewardPostHandler)
}

func (s *AccrualSystem) SetStorage(storage storage.StorageAccrualSystem) {
	s.storage = storage
	if err := s.storage.Open(); err != nil {
		slog.Error(err.Error())
		s.Shutdown()
	}

	// Try auto-migration
	s.storage.Migrate(s.Config.MigrDirName)

}

func (s *AccrualSystem) AddTask(task workeraccrual.Task) {
	slog.Debug("new task", slog.Any("order", task.Order))
	s.taskCh <- task
}
