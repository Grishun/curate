package domain

import (
	"context"
	"log/slog"
	"net/http"
)

type Storage interface {
	Get(ctx context.Context, currecny string, limit uint32) ([]Rate, error)
	GetAll(ctx context.Context, limit uint32) (map[string][]Rate, error)

	Insert(ctx context.Context, rates ...Rate) error

	HealthCheck(ctx context.Context) error
}

type Logger interface {
	With(args ...any) *slog.Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
}

type Provider interface {
	Name() string
	Fetch(ctx context.Context) (map[string]float64, error)
}

type HTTPClient interface {
	Do(ctx context.Context, opts ...RequestOption) (*http.Response, error)
}
