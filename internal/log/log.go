package log

import (
	"fmt"
	"log/slog"
)

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

// ----- wrappers -----

func (l *SlogProvider) Errorf(format string, v ...any) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l *SlogProvider) Warnf(format string, v ...any) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l *SlogProvider) Debugf(format string, v ...any) {
	l.Debug(fmt.Sprintf(format, v...))
}
