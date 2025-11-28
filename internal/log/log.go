package log

import "log/slog"

type SlogProvider struct {
	*slog.Logger
}

func NewSlog(opts ...Option) *SlogProvider {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	lg := slog.New(options.encoder)

	return &SlogProvider{lg}
}
