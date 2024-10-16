package servermart

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
	"github.com/mi4r/gophermart/lib/helper"

	"github.com/mi4r/gophermart/internal/auth"
)

const (
	successUserLogin   string = "user has been successfully registered and authenticated"
	orderAlreadyUpload string = "order number already uploaded by this user"
	orderAccepted      string = "order number accepted for processing"
	withdrawCompleted  string = "withdraw balance completed"
)

var (
	errEmptyLoginOrPassword = errors.New("login or password cannot be empty")
	errLoginIsExists        = errors.New("login already exists")
	errPasswordInvalid      = errors.New("invalid password")
	errUnauthorized         = errors.New("user is not authenticated")
	errInvalidOrderID       = errors.New("invalid order number format")
	errOrderUploadByAnother = errors.New("order number already uploaded by another user")
	errInsufficientFunds    = errors.New("insufficient funds")
)

// Ping
// @Description Простая проверка состояния сервера
// @Tags Разное
// @Success 200 {string} pong
// @Router /ping [get]
func (s *Gophermart) pingHandler(c echo.Context) error {
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
func (s *Gophermart) userRegisterHandler(c echo.Context) error {
	var creds storagemart.Creds
	if err := c.Bind(&creds); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if creds.IsEmpty() {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	user, err := storagemart.NewUserFromCreds(creds)
	if err != nil {
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Ожидаются еще ответы 409 - Логин уже занят
	if err := s.storage.UserCreate(context.Background(), user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.String(http.StatusConflict, errLoginIsExists.Error())
		}
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(auth.GetUserCookie(user.Login, s.Config.SecretKey))
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
func (s *Gophermart) userLoginHandler(c echo.Context) error {
	var creds storagemart.Creds
	if err := c.Bind(&creds); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if creds.IsEmpty() {
		return c.String(http.StatusBadRequest, errEmptyLoginOrPassword.Error())
	}

	user, err := s.storage.UserReadOne(context.Background(), creds.Login)
	if err == pgx.ErrNoRows {
		return c.String(http.StatusUnauthorized, pgx.ErrNoRows.Error())
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Если НЕ такой же то False
	if !user.PasswordCompare(creds) {
		return c.String(http.StatusUnauthorized, errPasswordInvalid.Error())
	}

	c.SetCookie(auth.GetUserCookie(user.Login, s.Config.SecretKey))
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
func (s *Gophermart) userPostOrdersHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c, s.Config.SecretKey)
	if !ok {
		return c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	bodyReader := c.Request().Body
	defer bodyReader.Close()

	bodyContent, err := io.ReadAll(bodyReader)
	if err != nil || len(bodyContent) == 0 {
		return c.String(http.StatusBadRequest, err.Error())
	}

	orderNumber := strings.TrimSpace(string(bodyContent))

	if !helper.IsLuhn(orderNumber) {
		return c.String(http.StatusUnprocessableEntity, errInvalidOrderID.Error())
	}

	// По идее тут отправляем в Accrual System
	//

	var emptyOrder storagemart.Order
	storedOrder, err := s.storage.UserOrderReadOne(context.Background(), orderNumber)
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
			// Тут ответ 200 по ТЗ
			return c.String(http.StatusOK, orderAlreadyUpload)
		}
	}

	err = s.storage.UserOrderCreate(context.Background(), login, orderNumber)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// Тут ответ 202 по ТЗ
	return c.String(http.StatusAccepted, orderAccepted)
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
func (s *Gophermart) userGetOrdersHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c, s.Config.SecretKey)
	if !ok {
		return c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	orders, err := s.storage.UserOrdersReadByLogin(context.Background(), login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(orders) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, orders)
}

// Balance get
// @Summary
// @Description Хендлер доступен только авторизованному пользователю.
// @Description В ответе должны содержаться данные о текущей сумме баллов лояльности,
// @Description а также сумме использованных за весь период регистрации баллов.
// @Tags Пользователь
// @Produce json
// @Success 200 {object} Balance "Успешная обработка запроса"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/balance [get]
func (s *Gophermart) userGetBalanceHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c, s.Config.SecretKey)
	if !ok {
		return c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	user, err := s.storage.UserReadOne(context.Background(), login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user.GetBalance())
}

// Balance withdraw
// @Summary
// @Description Хендлер доступен только авторизованному пользователю.
// @Description Номер заказа представляет собой гипотетический номер
// @Description нового заказа пользователя, в счёт оплаты которого списываются баллы.
// @Tags Заказы
// @Accept  json
// @Produce text/plain
// @Success 200 {string} string "Успешная обработка запроса"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 402 {string} string "На счету недостаточно средств"
// @Failure 422 {string} string "Неверный номер заказа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/balance/withdraw [post]
func (s *Gophermart) userBalanceWithdrawHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c, s.Config.SecretKey)
	if !ok {
		return c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}
	var req struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request format")
	}
	if !helper.IsLuhn(req.Order) {
		return c.String(http.StatusUnprocessableEntity, errInvalidOrderID.Error())
	}

	user, err := s.storage.UserReadOne(context.Background(), login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	curBalance := user.GetBalance().Current
	if curBalance < req.Sum {
		return c.String(http.StatusPaymentRequired, errInsufficientFunds.Error())
	}

	err = s.storage.WithdrawBalance(context.Background(), login, req.Order, req.Sum, curBalance)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, withdrawCompleted)
}

// Get balance withdrawals
// @Summary
// @Description Хендлер доступен только авторизованному пользователю.
// @Description Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым.
// @Description Формат даты — RFC3339.
// @Tags Заказы
// @Produce json
// @Success 200 {string} string "Успешная обработка запроса"
// @Failure 204 {string} string "Нет ни одного списания"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/withdrawals [get]
func (s *Gophermart) getBalanceWithdrawalsHandler(c echo.Context) error {
	login, ok := auth.ValidateUserCookie(c, s.Config.SecretKey)
	if !ok {
		return c.String(http.StatusUnauthorized, errUnauthorized.Error())
	}

	withdrawals, err := s.storage.GetUserWithdrawals(context.Background(), login)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(withdrawals) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, withdrawals)
}
