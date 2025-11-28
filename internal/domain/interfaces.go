package domain

import (
	"context"
	"log/slog"
)

type Storage interface {
	Get(ctx context.Context, code string) (Rate, error)
	GetAll(ctx context.Context) ([]Rate, error)

	Insert(ctx context.Context, rates ...Rate) error
}

type Logger interface {
	With(args ...any) *slog.Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Provider interface {
	Name() string
	Fetch(ctx context.Context) (map[string]float64, error)
}
