package serveraccrual

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	"github.com/mi4r/gophermart/lib/helper"
)

const (
	rewardCreated string = "new reward will be created"
	orderAccepted string = "order accepted"
)

var (
	errMatchKeyAlreadyExists     = errors.New("match key already exists")
	errInvalidReward             = errors.New("invalid reward format")
	errRewardIsNegative          = errors.New("reward value must not be a negative")
	errRewardIsInvalidType       = errors.New("reward type must be '%' or 'pt'")
	errInvalidRewardMatchIsEmpty = errors.New("match key must not be empty")
	errInvalidOrder              = errors.New("invalid order format")
	errNotFoundOrder             = errors.New("order not found")
	errInvalidOrderID            = errors.New("invalid order number format")
	errOrderAlreadyExists        = errors.New("order already exists")
)

// Ping
// @Description Простая проверка состояния сервера
// @Tags Разное
// @Success 200 {string} pong
// @Router /ping [get]
func (s *AccrualSystem) pingHandler(c echo.Context) error {
	storageOK := "pong"
	if err := s.storage.Ping(); err != nil {
		storageOK = err.Error()
	}
	return c.JSON(http.StatusOK, storageOK)
}

// Reward created
// @Summary Регистрация информации о вознаграждении за товар
// @Description Хендлер используется менеджерами для добавления механик вознаграждения за покупки
// @Description Полученные системой расчёта начислений составы чеков проверяются на совпадение с зарегистрированными в данном хендлере вознаграждениями
// @Tags Админ
// @Accept  application/json
// @Produce text/plain
// @Param reward body Reward true "Механика вознаграждения"
// @Success 200 {string} string "Вознаграждение успешно зарегистрировано"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 409 {string} string "Ключ поиска уже зарегистрирован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/goods [post]
func (s *AccrualSystem) rewardPostHandler(c echo.Context) error {
	var reward storageaccrual.Reward
	if err := c.Bind(&reward); err != nil {
		return c.JSON(http.StatusBadRequest, errInvalidReward.Error())
	}
	if reward.IsEmptyMatch() {
		return c.JSON(http.StatusBadRequest, errInvalidRewardMatchIsEmpty.Error())
	}

	if reward.IsNegative() {
		return c.JSON(http.StatusBadRequest, errRewardIsNegative.Error())
	}

	if !reward.IsValidType() {
		return c.JSON(http.StatusBadRequest, errRewardIsInvalidType.Error())
	}

	if err := s.storage.RewardCreate(reward); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.String(http.StatusConflict, errMatchKeyAlreadyExists.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, rewardCreated)
}

// Registers a new order
// @Summary Регистрация нового совершённого заказа
// @Description Для начисления баллов состав заказа должен быть проверен на совпадения с зарегистрированными записями вознаграждений за товары
// @Description Начисляется сумма совпадений
// @Description Принятый заказ не обязан браться в обработку непосредственно в момент получения запроса
// @Tags Админ
// @Accept  application/json
// @Produce text/plain
// @Param reward body Order true "Регистрация нового совершённого заказа"
// @Success 202 {string} string "Заказ успешно принят в обработку"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 409 {string} string "Заказ уже принят в обработку"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/orders [post]
func (s *AccrualSystem) ordersPostHandler(c echo.Context) error {
	var order storageaccrual.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, errInvalidOrder.Error())
	}

	// Нужна ли проверка? поидее этой ручкой наполняем базу и уже проверили все
	if !helper.IsLuhn(order.Order) {
		return c.String(http.StatusBadRequest, errInvalidOrderID.Error())
	}

	if err := s.storage.OrderRegCreate(order); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			slog.Debug("internal error. 23505", slog.String("msg", err.Error()))
			return c.String(http.StatusConflict, errOrderAlreadyExists.Error())
		}
		slog.Debug("internal error. unknown", slog.String("msg", err.Error()))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusAccepted, orderAccepted)
}

// Get accrual info
// @Summary Получение информации о расчёте начислений
// @Description Получение информации о расчёте начислений баллов лояльности за совершённый заказ
// @Description Номером заказа является последовательность цифр произвольной длины.
// @Description Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
// @Tags Сервис
// @Accept text/plain
// @Produce  application/json
// @Param number path string true "Номером заказа"
// @Success 200 {object} storagedefault.Order "Успешная обработка запроса"
// @Success 204 {string} string "Заказ не зарегистрирован в системе расчёта"
// @Failure 429 {string} string "Превышено количество запросов к сервису"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/orders/{number} [get]
func (s *AccrualSystem) ordersGetHandler(c echo.Context) error {
	number := c.Param("number")
	if !helper.IsLuhn(number) {
		// Нет более подходящего статуса ответа исходя из ТЗ
		return c.String(http.StatusNoContent, errInvalidOrderID.Error())
	}

	// TODO 429 status
	// ...
	// ctx, cancel := context.WithCancel(context.TODO())
	ctx := context.Background()
	order, err := s.storage.OrderRegReadOne(ctx, number)
	if err != nil {
		// TODO. Обработка ошибок, пока будет просто ненаход
		return c.String(http.StatusNoContent, errNotFoundOrder.Error())

		// TODO. 500 status
		// ...
	}
	return c.JSON(http.StatusOK, order)
}
