package middleware

import (
	"context"
	"strconv"
	"test_task/internal/dto"
	"test_task/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		// Генерация request ID
		guid := uuid.New()
		reqID := strconv.FormatUint(uint64(guid.ID()), 10)

		// Получение оригинального контекста
		ctx := c.Request().Context()

		// Создание нового контекста с request ID
		ctx = context.WithValue(ctx, dto.RequestID, reqID)

		// Обновление контекста в echo.Context
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		// Теперь можешь использовать этот контекст
		logger.GetLoggerFromCtx(ctx).Debug(ctx,
			"request", zap.String("method", c.Request().Method),
			zap.Time("request_time", time.Now()), zap.String("request_id", reqID))

		// Передача дальше
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}

}
