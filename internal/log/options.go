package log

import (
	"log/slog"
	"os"
)

type Option func(l Options)

type Options struct {
	encoder slog.Handler
}

func NewOptions() Options {
	return Options{
		encoder: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}
}

func WithEncoderJSON(level slog.Level) Option {
	return func(opt Options) {
		opt.encoder = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
}

func WithEncoderText(level slog.Level) Option {
	return func(opt Options) {
		opt.encoder = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
}
