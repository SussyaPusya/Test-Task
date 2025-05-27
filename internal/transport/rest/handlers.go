package rest

import (
	"context"
	"net/http"
	"test_task/internal/dto"
	"test_task/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Service interface {
	AddPeople(ctx context.Context, person *dto.Person) error
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) Ping(c echo.Context) error {

	return c.JSON(http.StatusOK, "PONG!")
}

func (h *Handlers) AddPeople(c echo.Context) error {
	ctx := c.Request().Context()

	var reqStruct dto.Person

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Trying to bind JSON to dto.Person")

	if err := c.Bind(&reqStruct); err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to bind JSON: "+err.Error())
		return c.JSON(http.StatusBadRequest, "bad json")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Received request to add person: Name="+reqStruct.Name+", Surname="+reqStruct.Surname)

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Calling service.AddPeople")

	err := h.service.AddPeople(ctx, &reqStruct)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Service failed to add person: "+err.Error())
		return c.JSON(http.StatusInternalServerError, "InternalServerError")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person successfully added to database")
	return c.JSON(http.StatusOK, "add Person")
}
