package logger

import (
	"context"
	"test_task/internal/config"
	"test_task/internal/dto"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	L *zap.Logger
}

func New(ctx context.Context, config *config.Logger) (context.Context, error) {
	var level = zapcore.InfoLevel
	if config.Level == "DEBUG" {
		level = zap.DebugLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, dto.Logger, &Logger{logger})

	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(dto.Logger).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(dto.RequestID) != nil {
		fields = append(fields, zap.String(string(dto.RequestID), ctx.Value(dto.RequestID).(string)))
	}

	l.L.Info(msg, fields...)
}
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(dto.RequestID) != nil {
		fields = append(fields, zap.String(string(dto.RequestID), ctx.Value(dto.RequestID).(string)))
	}

	l.L.Debug(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(dto.RequestID) != nil {
		fields = append(fields, zap.String(string(dto.RequestID), ctx.Value(dto.RequestID).(string)))
	}

	l.L.Fatal(msg, fields...)
}
