package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

// Ping
// @Summary Health check of the server
// @Tags Common
// @Accept  text/plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (s *Server) pingHandler(c echo.Context) error {
	return c.String(200, "pong")
}

func (s *Server) registerHandler(c echo.Context) error {
	var user storage.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	if strings.TrimSpace(user.Login) == "" || strings.TrimSpace(user.Password) == "" {
		return c.JSON(http.StatusBadRequest, "Empty login or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Server can't hash the password")
	}

	_, err = s.storage.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", user.Login, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.JSON(http.StatusConflict, "Login already taken")
		}
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	auth.SetUserCookie(c, user.Login)
	return nil
}

func (s *Server) loginHandler(c echo.Context) error {
	return nil
}
