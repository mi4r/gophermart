package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"

	"github.com/mi4r/gophermart/internal/auth"
)

// const (
// 	errInvalidRequest      string = "invalid request"
// 	errLoginOrPassEmpty    string = "login or password cannot be empty"
// 	errHashedPasswordError string = "server can't hash the password"
// 	errUserNotFound        string = "user not found"
// )

var errEmptyLoginOrPassword = errors.New("login or password cannot be empty")

type Resp struct {
	Text    string `json:"text"`
	Payload any    `json:"payload"`
} //@name Response

func newErrResp(err error) Resp {
	return Resp{
		Text: err.Error(),
	}
}

// Ping
// @Summary Health check of the server
// @Tags Common
// @Success 200 {object} Response
// @Router /ping [get]
func (s *Server) pingHandler(c echo.Context) error {
	return c.JSON(200, Resp{
		Text: "pong",
	})
}

// User register
// @Summary Check creds and registration user
// @Tags users
// @Accept  json
// @Param creds body Creds true "login and password"
// @Router /api/user/register [post]
// @Success 200
// @failure 400 {object} Response
// @failure 409 {object} Response
// @failure 500 {object} Response
func (s *Server) registerHandler(c echo.Context) error {
	var user storage.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp(err))
	}

	if user.Login == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, newErrResp(errEmptyLoginOrPassword))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp(err))
	}
	user.Password = string(hashedPassword)

	// Ожидаются еще ответы 409 - Логин уже занят
	if err := s.storage.UserCreate(user); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp(err))
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
// @failure 400 {object} Response
// @failure 401 {object} Response
// @failure 500 {object} Response
// @Router /api/user/login [post]
func (s *Server) loginHandler(c echo.Context) error {
	return nil
}
