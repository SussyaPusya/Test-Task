package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"test_task/internal/dto"
	"test_task/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Service interface {
	AddPeople(ctx context.Context, person *dto.Person) (string, error)
	DeletePerson(ctx context.Context, id string) error
	GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error)
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

	id, err := h.service.AddPeople(ctx, &reqStruct)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Service failed to add person: "+err.Error())
		return c.JSON(http.StatusInternalServerError, "InternalServerError")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person successfully added to database")
	return c.JSON(http.StatusOK, map[string]string{"status": "succesful", "ID": id})
}
func (h *Handlers) DeletePerson(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.QueryParam("id")

	err := h.service.DeletePerson(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "bla bla")
	}

	return c.JSON(http.StatusOK, "succesful delete")
}

func (h *Handlers) GetPeople(c echo.Context) error {
	ctx := c.Request().Context()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	filter := dto.PersonFilter{
		Name:     c.QueryParam("name"),
		Surname:  c.QueryParam("surname"),
		Age:      c.QueryParam("age"),
		Patronym: c.QueryParam("patronymic"),
		Gender:   c.QueryParam("gender"),
		Country:  c.QueryParam("country"),
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx, fmt.Sprintf("Fetching people with filter: %+v, limit=%d, offset=%d", filter, limit, offset))

	people, err := h.service.GetPeople(ctx, &filter, limit, offset)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to fetch people: "+err.Error())
		return c.JSON(http.StatusInternalServerError, "internal error")
	}


	

	return c.JSON(http.StatusOK, people)
}
