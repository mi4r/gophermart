package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

var errEmptyLoginOrPassword = errors.New("login or password cannot be empty")

// Ping
// @Summary Health check of the server
// @Tags Common
// @Success 200 {string} pong
// @Router /ping [get]
func (s *Server) pingHandler(c echo.Context) error {
	return c.String(200, "pong")
}

// User register
// @Summary Check creds and registration user
// @Tags users
// @Accept  json
// @Param creds body Creds true "login and password"
// @Router /api/user/register [post]
// @Success 200
// @Failure 400 {string}
// @Failure 409 {string}
// @Failure 500 {string}
func (s *Server) registerHandler(c echo.Context) error {
	var user storage.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if user.Login == "" || user.Password == "" {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(hashedPassword)

	// Ожидаются еще ответы 409 - Логин уже занят
	if err := s.storage.UserCreate(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	auth.SetUserCookie(c, user.Login)
	return nil
}

// User login
// @Summary Check creds and login user
// @Tags users
// @Accept  json
// @Produce json
// @Param creds body Creds true "login and password"
// @Success 200
// @Failure 400 {string}
// @Failure 401 {string}
// @Failure 500 {string}
// @Router /api/user/login [post]
func (s *Server) loginHandler(c echo.Context) error {
	return nil
}
