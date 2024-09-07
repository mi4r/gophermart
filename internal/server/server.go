package server

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/storage"
)

type Server struct {
	config  config.ServerConfig
	storage *storage.Storage
	router  *echo.Echo
}

func NewServer(config config.ServerConfig, storage *storage.Storage) *Server {
	return &Server{
		config:  config,
		storage: storage,
		router:  echo.New(),
	}
}

func (s *Server) Start() {
	s.Configure()
	slog.Error(s.router.Start(s.config.ListenAddr).Error())
	os.Exit(1)
}

func (s *Server) Configure() {
	s.setMiddlewares()
	s.setRoutes()
	// ...
}

func (s *Server) setRoutes() {
	s.router.GET("/ping", s.pingHandler)
	// ...
}

func (s *Server) setMiddlewares() {
	s.router.Use(middleware.Logger())
	// ...
}
