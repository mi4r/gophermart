package serveraccrual

import (
	"log/slog"

	"github.com/mi4r/gophermart/internal/server"
	"github.com/mi4r/gophermart/internal/storage"
)

type AccrualSystem struct {
	*server.Server
	storage storage.StorageAccrualSystem
}

func NewAccrualSystem(server *server.Server) *AccrualSystem {
	return &AccrualSystem{
		Server: server,
	}
}

func (s *AccrualSystem) SetRoutes() {
	s.Router.GET("/ping", s.pingHandler)
	// TODO
	gAPI := s.Router.Group("/api")
	gAPI.GET("/orders/:number", s.ordersGetHandler)
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
