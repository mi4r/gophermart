package server

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/mi4r/gophermart/internal/storage"

	"github.com/mi4r/gophermart/internal/auth"
)

const (
	successUserLogin   string = "user has been successfully registered and authenticated"
	orderAlreadyUpload string = "order number already uploaded by this user"
	orderAccepted      string = "order number accepted for processing"
)

var (
	errEmptyLoginOrPassword = errors.New("login or password cannot be empty")
	errLoginIsExists        = errors.New("login already exists")
	errPasswordInvalid      = errors.New("invalid password")
	errUnauthorized         = errors.New("user is not authenticated")
	errInvalidOrderID       = errors.New("invalid order number format")
	errOrderUploadByAnother = errors.New("order number already uploaded by another user")
)

// Ping
// @Description Простая проверка состояния сервера
// @Tags Разное
// @Success 200 {string} pong
// @Router /ping [get]
func (s *Server) pingHandler(c echo.Context) error {
	storageOK := "pong"
	if err := s.storage.Ping(); err != nil {
		storageOK = err.Error()
	}
	return c.JSON(http.StatusOK, storageOK)
}

// User register
// @Summary Регистрация пользователя
// @Description Для передачи аутентификационных данных используется механизм cookies
// @Tags Пользователь
// @Accept  json
// @Param creds body Creds true "Логин и пароль не зарегистрированного пользователя"
// @Router /api/user/register [post]
// @Success 200 {string} string "Пользователь успешно зарегистрирован и аутентифицирован"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 409 {string} string "Логин уже занят"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
func (s *Server) userRegisterHandler(c echo.Context) error {
	var creds storage.Creds
	if err := c.Bind(&creds); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if creds.IsEmpty() {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	user, err := storage.NewUserFromCreds(creds)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Ожидаются еще ответы 409 - Логин уже занят
	if err := s.storage.UserCreate(user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.String(http.StatusConflict, errLoginIsExists.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(auth.GetUserCookie(user.Login))
	return c.String(http.StatusOK, successUserLogin)
}

// User login
// @Summary Аутентификация пользователя
// @Description Для передачи аутентификационных данных используется механизм cookies
// @Tags Пользователь
// @Accept  json
// @Produce text/plain
// @Param creds body Creds true "Логин и пароль зарегистрированного пользователя"
// @Success 200 {string} string "Пользователь успешно зарегистрирован и аутентифицирован"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Неверная пара логин/пароль"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/login [post]
func (s *Server) userLoginHandler(c echo.Context) error {
	var creds storage.Creds
	if err := c.Bind(&creds); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if creds.IsEmpty() {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	user, err := s.storage.UserReadOne(creds.Login)
	if err == pgx.ErrNoRows {
		return c.String(http.StatusUnauthorized, pgx.ErrNoRows.Error())
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Если НЕ такой же то False
	if !user.PasswordCompare(creds) {
		return c.String(http.StatusUnauthorized, errPasswordInvalid.Error())
	}

	c.SetCookie(auth.GetUserCookie(user.Login))
	return c.String(http.StatusOK, successUserLogin)
}

// Order register
// @Summary Загрузка номера заказа
// @Description Хендлер доступен только аутентифицированным пользователям
// @Description Номером заказа является последовательность цифр произвольной длины
// @Tags Заказы
// @Accept  text/plain
// @Produce text/plain
// @Param number body string true "Трек номер заказа"
// @Success 200 {string} string "Номер заказа уже был загружен этим пользователем"
// @Success 202 {string} string "Новый номер заказа принят в обработку"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 409 {string} string "Номер заказа уже был загружен другим пользователем"
// @Failure 422 {string} string "Неверный формат номера заказа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/orders [post]
func (s *Server) userPostOrdersHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c)
	if !ok {
		c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	bodyReader := c.Request().Body
	defer bodyReader.Close()

	bodyContent, err := io.ReadAll(bodyReader)
	if err != nil || len(bodyContent) == 0 {
		return c.String(http.StatusBadRequest, err.Error())
	}

	orderNumber := strings.TrimSpace(string(bodyContent))
	if !checkLuhn(orderNumber) {
		return c.String(http.StatusUnprocessableEntity, errInvalidOrderID.Error())
	}
	var emptyOrder storage.Order
	storedOrder, err := s.storage.OrderReadOne(orderNumber)
	if err != nil && err != pgx.ErrNoRows {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Если заказ был найден
	if storedOrder != emptyOrder {
		// Проверка на то что заказ уже не создан другим пользователем
		if storedOrder.UserLogin != login {
			return c.String(http.StatusConflict, errOrderUploadByAnother.Error())
			// Проверка на то что заказ уже не был создан этим же пользователем
		} else if storedOrder.UserLogin == login {
			return c.String(http.StatusOK, orderAlreadyUpload)
		}
	}

	err = s.storage.OrderCreate(login, orderNumber)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, orderAccepted)
}

func checkLuhn(orderID string) bool {
	sum := 0
	nDigits := len(orderID)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit, err := strconv.Atoi(string(orderID[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (sum % 10) == 0
}

// Orders get
// @Summary Получение списка загруженных номеров заказов
// @Description Хендлер доступен только авторизованному пользователю
// @Description Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым
// @Description Формат даты — RFC3339.
// @Tags Заказы
// @Accept  text/plain
// @Produce json
// @Success 200 {object} []Order "Успешная обработка запроса"
// @Success 204 {string} string "Нет данных для ответа"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/orders [get]
func (s *Server) userGetOrdersHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c)
	if !ok {
		c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	orders, err := s.storage.OrdersReadByLogin(login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(orders) == 0 {
		c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, orders)
}
