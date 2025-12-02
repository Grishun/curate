package influx

import (
	"context"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type Options struct {
	ctx          context.Context
	writeTimeout time.Duration
	queryTimeout time.Duration
	logger       domain.Logger
	currencies   []string
}

type Option func(*Options)

func WithLogger(l domain.Logger) Option {
	return func(opt *Options) {
		opt.logger = l
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(opt *Options) {
		opt.writeTimeout = timeout
	}
}

func WithQueryTimeout(timeout time.Duration) Option {
	return func(opt *Options) {
		opt.queryTimeout = timeout
	}
}

func WithContext(ctx context.Context) Option {
	return func(opt *Options) {
		opt.ctx = ctx
	}
}

func WithCurrencies(currencies ...string) Option {
	return func(opt *Options) {
		opt.currencies = currencies
	}
}

// NewOptions returns an empty Options struct! REQUIRED to fill it with options
func NewOptions(opts ...Option) *Options {
	return &Options{
		logger:       log.NewSlog(),
		writeTimeout: 10 * time.Second,
		queryTimeout: time.Minute,
		ctx:          context.Background(),
		currencies:   []string{"BTC", "ETH", "TRX"},
	}
}
