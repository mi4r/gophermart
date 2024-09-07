package server

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/storage"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	if err := s.router.Start(s.config.ListenAddr); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (s *Server) Shutdown() {
	slog.Debug("http server stopped")
	// ...
}

func (s *Server) Configure() {
	s.setMiddlewares()
	s.setRoutes()
	// ...
}

func (s *Server) setRoutes() {
	// swagger
	s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	s.router.GET("/ping", s.pingHandler)

	// ...
}

func (s *Server) setMiddlewares() {
	s.router.Use(middleware.Logger())

}
