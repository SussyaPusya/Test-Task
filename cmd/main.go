package main

import (
	"context"
	"test_task/internal/clients"
	"test_task/internal/config"
	"test_task/internal/repository"
	"test_task/internal/service"
	"test_task/internal/transport/rest"
	"test_task/internal/transport/rest/middleware"
	"test_task/pkg/logger"
	"test_task/pkg/postgres"
)

func main() {
	ctx := context.Background()

	config, _ := config.NewConfig()
	ctx, _ = logger.New(ctx, &config.Logger)

	cleints := clients.NewExtanlAPI(&config.ExternalAPI)

	pg, err := postgres.NewPostgres(ctx, &config.Postgres)
	if err != nil {
		return
	}

	repo := repository.NewRepository(pg)

	service := service.NewService(cleints, repo)

	handle := rest.NewHandlers(service)

	midle := middleware.NewMiddleware()

	router := rest.NewRouter(&config.Rest, handle, ctx, midle)

	router.Run(ctx)
}
