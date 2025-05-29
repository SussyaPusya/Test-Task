package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"test_task/internal/config"
	"test_task/internal/transport/rest/middleware"
	"test_task/pkg/logger"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Router struct {
	router *echo.Echo

	handlers *Handlers

	config *config.Rest
}

func NewRouter(cfg *config.Rest, handlers *Handlers, ctx context.Context, middleware *middleware.Middleware) *Router {
	e := echo.New()

	e.Server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	e.Use(middleware.Logger)
	e.GET("/ping", handlers.Ping)
	e.GET("/people", handlers.GetPeople)
	e.POST("/people", handlers.AddPeople)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.DELETE("/people/delete", handlers.DeletePerson)
	e.PATCH("/people/update", handlers.UpdatePerson)

	return &Router{router: e, handlers: handlers, config: cfg}

}

func (r *Router) Run(ctx context.Context) {

	restAddr := fmt.Sprintf(":%d", r.config.Port)

	if err := r.router.Start(restAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to start server", zap.Error(err))
	}

}
