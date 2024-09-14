package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mi4r/gophermart/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage — это мок, который реализует интерфейс Storage.
type MockStorage struct {
	mock.Mock
}

// UserCreate — мок реализация метода для создания пользователя.
func (m *MockStorage) UserCreate(user storage.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// UserReadOne — мок реализация метода для чтения данных пользователя.
func (m *MockStorage) UserReadOne(login string) (storage.User, error) {
	args := m.Called(login)
	return args.Get(0).(storage.User), args.Error(1)
}

func (m *MockStorage) UserReadAll() ([]storage.User, error) {
	args := m.Called()
	return args.Get(0).([]storage.User), args.Error(1)
}

func (m *MockStorage) Open() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Close() {
	return
}

func TestUserRegisterHandler(t *testing.T) {
	e := echo.New()
	mockStorage := new(MockStorage)
	server := &Server{storage: mockStorage}

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
