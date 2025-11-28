package domain

import (
	"context"
	"log/slog"
)

type Storage interface {
	Get(ctx context.Context, code string) (domain.Rate, error)
	GetAll(ctx context.Context) ([]domain.Rate, error)

	Insert(ctx context.Context, rates ...domain.Rate) error
}

type Logger interface {
	With(args ...any) *slog.Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
