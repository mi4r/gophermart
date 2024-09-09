package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

const (
	InvalidRequest      string = "Invalid request"
	HashedPasswordError string = "Server can't hash the password"
	UserNotFound        string = "User not found"
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
		return c.JSON(http.StatusBadRequest, InvalidRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HashedPasswordError)
	}
	user.Password = string(hashedPassword)

	if err := s.storage.UserCreate(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	auth.SetUserCookie(c, user.Login)
	return nil
}

func (s *Server) loginHandler(c echo.Context) error {
	return nil
}
