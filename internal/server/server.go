package server

import (
	"context"
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
	storage storage.Storage
	router  *echo.Echo
}

func NewServer(config config.ServerConfig, storage storage.Storage) *Server {
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
	s.router.POST("/api/user/register", s.registerHandler)
	s.router.POST("/api/user/login", s.loginHandler)
	// ...
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
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
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
