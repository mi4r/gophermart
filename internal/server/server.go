package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mi4r/gophermart/internal/storage"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	GophermartName    = "GOPHERMART"
	AccrualSystemName = "ACCRUAL_SYS"
)

type Config struct {
	ServiceName string
	Listen      string
	SecretKey   string
}

type Server struct {
	config  Config
	storage storage.Storage
	router  *echo.Echo
}

func NewServer(config Config, storage storage.Storage) *Server {
	return &Server{
		config:  config,
		storage: storage,
		router:  echo.New(),
	}
}

func (s *Server) Start() {
	s.Configure()
	if err := s.router.Start(s.config.Listen); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (s *Server) Shutdown() {
	slog.Debug("http server stopped")
	// ...
	os.Exit(0)
}

func (s *Server) Configure() {
	s.setMiddlewares()
	s.setDefaultRoutes()
	s.setStorage()
	// ...
}

// TODO
func (s *Server) setStorage() {
	if err := s.storage.Open(); err != nil {
		slog.Error(err.Error())
		s.Shutdown()
	}

}

func (s *Server) setDefaultRoutes() {
	// swagger
	s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	s.router.GET("/ping", s.pingHandler)
}

func (s *Server) setMiddlewares() {
	s.setLogger()
	// s.router.Use(middleware.Logger())

}

func (s *Server) setLogger() {
	s.router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("service", s.config.ServiceName),
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("service", s.config.ServiceName),
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
