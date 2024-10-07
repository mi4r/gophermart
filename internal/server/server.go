package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	GophermartName = "GOPHERMART"
	AccrualName    = "ACCRUAL"
)

type Config struct {
	ServiceName string
	Listen      string
	SecretKey   string
	MigrDirName string
}

type Server struct {
	Config Config
	Router *echo.Echo
}

func NewServer(Config Config) *Server {
	return &Server{
		Config: Config,
		Router: echo.New(),
	}
}

func (s *Server) Start() {
	s.Configure()
	if err := s.Router.Start(s.Config.Listen); err != nil {
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
	// s.setStorage()
	// ...
}

// TODO

func (s *Server) setDefaultRoutes() {
	// swagger
	s.Router.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Server) setMiddlewares() {
	s.setLogger()
	// s.Router.Use(middleware.Logger())

}

func (s *Server) setLogger() {
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("service", s.Config.ServiceName),
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("service", s.Config.ServiceName),
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
