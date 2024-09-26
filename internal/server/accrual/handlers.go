package serveraccrual

import (
	"net/http"

	"github.com/labstack/echo/v4"
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
)

const (
	goodsCreated string = "new goods will be created"
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
	var _ storageaccrual.Reward
	return c.String(http.StatusOK, goodsCreated)
}
