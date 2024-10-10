package servermart

import (
	"log/slog"

	"github.com/mi4r/gophermart/internal/server"
	"github.com/mi4r/gophermart/internal/storage"
)

type Gophermart struct {
	*server.Server
	storage storage.StorageGophermart
}

func NewGophermart(server *server.Server) *Gophermart {
	return &Gophermart{
		Server: server,
	}
}

func (s *Gophermart) SetStorage(storage storage.StorageGophermart) {
	s.storage = storage
	if err := s.storage.Open(); err != nil {
		slog.Error(err.Error())
		s.Shutdown()
	}

	// Try auto-migration
	s.storage.Migrate(s.Config.MigrDirName)

}

func (s *Gophermart) SetRoutes() {
	s.Router.GET("/ping", s.pingHandler)
	gUsers := s.Router.Group("/api/user")
	gUsers.POST("/register", s.userRegisterHandler)
	gUsers.POST("/login", s.userLoginHandler)
	gUsers.POST("/orders", s.userPostOrdersHandler)
	gUsers.GET("/orders", s.userGetOrdersHandler)
	gUsers.GET("/balance", s.userGetBalanceHandler)
	gUsers.POST("/balance/withdraw", s.userBalanceWithdrawHandler)
	gUsers.GET("/withdrawals", s.getBalanceWithdrawalsHandler)
}
