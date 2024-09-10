package server

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

const (
	InvalidRequest      string = "Invalid request"
	HashedPasswordError string = "Server can't hash the password"
	UserNotFoundError   string = "User not found"
	LoginUsedError      string = "Login has already taken"
	InvalidPassword     string = "Invalid password"
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
		return c.JSON(http.StatusConflict, err.Error())
	}

	cookie := auth.GetUserCookie(user.Login)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "The user has been successfully registered and authenticated")
}

func (s *Server) loginHandler(c echo.Context) error {
	var user storage.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidRequest)
	}

	userStored, err := s.storage.UserReadOne(user.Login)
	if err == pgx.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, UserNotFoundError)
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userStored.Password), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, InvalidPassword)
	}

	cookie := auth.GetUserCookie(user.Login)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "The user has been successfully authenticated")
}
