package server

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

var errEmptyLoginOrPassword = errors.New("login or password cannot be empty")

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

	if err := s.storage.UserCreate(user); err != nil {
		return c.String(http.StatusConflict, err.Error())
	}

	cookie := auth.GetUserCookie(user.Login)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "The user has been successfully registered and authenticated")
}

func (s *Server) loginHandler(c echo.Context) error {
	var user storage.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if user.Login == "" || user.Password == "" {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	userStored, err := s.storage.UserReadOne(user.Login)
	if err == pgx.ErrNoRows {
		return c.String(http.StatusUnauthorized, err.Error())
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userStored.Password), []byte(user.Password)); err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	cookie := auth.GetUserCookie(user.Login)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "The user has been successfully authenticated")
}
