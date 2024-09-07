package server

import "github.com/labstack/echo/v4"

// Ping
// @Summary Health check of the server
// @Tags Common
// @Accept  text/plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (s *Server) pingHandler(c echo.Context) error {
	return c.String(200, "pong")
}
