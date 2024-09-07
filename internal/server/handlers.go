package server

import "github.com/labstack/echo/v4"

func (s *Server) pingHandler(c echo.Context) error {
	return c.String(200, "pong")
}
