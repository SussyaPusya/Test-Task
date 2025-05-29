package main

import (
	"context"
	_ "test_task/docs"
	"test_task/internal/clients"
	"test_task/internal/config"
	"test_task/internal/repository"
	"test_task/internal/service"
	"test_task/internal/transport/rest"
	"test_task/internal/transport/rest/middleware"
	"test_task/pkg/logger"
	"test_task/pkg/postgres"
)

// @title People Enrichment API
// @version 1.0
// @description API для добавления и фильтрации людей
// @host localhost:8080
// @BasePath /

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
