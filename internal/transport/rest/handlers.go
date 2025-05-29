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
	UpdatePerson(ctx context.Context, person *dto.Person) error
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

// AddPeople godoc
// @Summary Добавить нового человека
// @Description Считывает с внешнего апи пол возраст и гендер
// @Tags people
// @Accept  json
// @Produce  json
// @Param   person  body dto.Person true "Person info"
// @Success 200 {string} string "ID"
// @Failure 400 {string} string "bad json"
// @Failure 500 {string} string "internal error"
// @Router /people [post]
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

// DeletePerson godoc
// @Summary Удалить человека по ID
// @Description Удаляет запись из базы данных по идентификатору
// @Tags people
// @Accept  json
// @Produce  json
// @Param   id   query string true "ID человека"
// @Success 200 {string} string "succesful delete"
// @Failure 500 {string} string "bla bla"
// @Router /people/delete [delete]
func (h *Handlers) DeletePerson(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.QueryParam("id")

	err := h.service.DeletePerson(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "bla bla")
	}

	return c.JSON(http.StatusOK, "succesful delete")
}

// GetPeople godoc
// @Summary Получить список людей
// @Description Получает список людей с фильтрацией и пагинацией
// @Tags people
// @Accept  json
// @Produce  json
// @Param   name        query string false "Имя"
// @Param   surname     query string false "Фамилия"
// @Param   patronymic  query string false "Отчество"
// @Param   gender      query string false "Пол"
// @Param   age         query string false "Возраст"
// @Param   country     query string false "Национальность"
// @Param   page        query int false "Номер страницы (по умолчанию 1)"
// @Param   limit       query int false "Размер страницы (по умолчанию 10)"
// @Success 200 {array} dto.Person
// @Failure 500 {string} string "internal error"
// @Router /people [get]
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

// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновляет существующую запись по ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param   id      query string true "ID человека"
// @Param   person  body dto.Person true "Обновлённые данные человека"
// @Success 200 {string} string "succesful update"
// @Failure 400 {string} string "bad json"
// @Failure 500 {string} string "InternalServerError"
// @Router /people/update [patch]
func (h *Handlers) UpdatePerson(c echo.Context) error {
	ctx := c.Request().Context()

	var reqStruct dto.Person
	id := c.QueryParam("id")

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Trying to bind JSON to dto.Person")

	if err := c.Bind(&reqStruct); err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to bind JSON: "+err.Error())
		return c.JSON(http.StatusBadRequest, "bad json")
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Received request to add person: Name="+reqStruct.Name+", Surname="+reqStruct.Surname)

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Calling service.UpdatePerson")
	reqStruct.ID = id

	err := h.service.UpdatePerson(ctx, &reqStruct)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Service failed to update person: "+err.Error())
		return c.JSON(http.StatusInternalServerError, "InternalServerError")
	}

	return c.JSON(http.StatusOK, "succesful update")
}
