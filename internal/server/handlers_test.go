package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mi4r/gophermart/internal/config"
	"github.com/mi4r/gophermart/internal/mocks"
	"github.com/mi4r/gophermart/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRegisterHandler(t *testing.T) {
	// Примерно так делается
	ctrl := gomock.NewController(t)
	m := mocks.NewMockStorage(ctrl)
	c := config.ServerConfig{}
	server := NewServer(c, m)
	tests := []struct {
		name           string
		requestBody    string
		mockUserCreate error
		expectedCode   int
		expectedBody   string
	}{
		{
			name:           "Successful Registration",
			requestBody:    `{"login": "testuser", "password": "testpass"}`,
			mockUserCreate: nil,
			expectedCode:   http.StatusOK,
			expectedBody:   successUserLogin,
		},
		{
			name:           "Login Already Exists",
			requestBody:    `{"login": "testuser", "password": "testpass"}`,
			mockUserCreate: &pgconn.PgError{Code: "23505"},
			expectedCode:   http.StatusConflict,
			expectedBody:   errLoginIsExists.Error(),
		},
		{
			name:           "Empty Login or Password",
			requestBody:    `{"login": "", "password": ""}`,
			mockUserCreate: nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   errEmptyLoginOrPassword.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On("UserCreate", mock.Anything).Return(tt.mockUserCreate)

			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, server.userRegisterHandler(c)) {
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestUserLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockStorage(ctrl)
	c := config.ServerConfig{}
	server := NewServer(c, m)

	tests := []struct {
		name            string
		requestBody     string
		mockUser        storage.User
		mockUserError   error
		mockPasswordCmp bool
		expectedCode    int
		expectedBody    string
	}{
		{
			name:            "Successful Login",
			requestBody:     `{"login": "testuser", "password": "testpass"}`,
			mockUser:        storage.User{Creds: storage.Creds{Login: "testuser", Password: "hashedpass"}},
			mockUserError:   nil,
			mockPasswordCmp: true,
			expectedCode:    http.StatusOK,
			expectedBody:    successUserLogin,
		},
		{
			name:            "Invalid Login or Password",
			requestBody:     `{"login": "testuser", "password": "wrongpass"}`,
			mockUser:        storage.User{Creds: storage.Creds{Login: "testuser", Password: "hashedpass"}},
			mockUserError:   nil,
			mockPasswordCmp: false,
			expectedCode:    http.StatusUnauthorized,
			expectedBody:    errPasswordInvalid.Error(),
		},
		{
			name:            "User Not Found",
			requestBody:     `{"login": "nonexistent", "password": "testpass"}`,
			mockUser:        storage.User{},
			mockUserError:   pgx.ErrNoRows,
			mockPasswordCmp: false,
			expectedCode:    http.StatusUnauthorized,
			expectedBody:    pgx.ErrNoRows.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On("UserReadOne", mock.Anything).Return(tt.mockUser, tt.mockUserError)
			mockStorage.On("PasswordCompare", mock.Anything).Return(tt.mockPasswordCmp)

			req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, server.userLoginHandler(c)) {
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}
