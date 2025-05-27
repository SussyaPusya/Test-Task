package postgres

import (
	"context"
	"errors"
	"fmt"
	"test_task/internal/config"
	"test_task/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPostgres(ctx context.Context, config *config.Postgres) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.MaxConn,
		config.MinConn,
	)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "unable to connect to database:", zap.Error(err))
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Database is not connected", zap.Error(err))
	}

	migraton, err := migrate.New(
		"file://././db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			config.User,
			config.Password,
			config.Host,
			config.Port,

			config.Database,
		))
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	if err := migraton.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "unable to run migrations:", zap.Error(err))
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return conn, nil
}
