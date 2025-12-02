package influx

import (
	"context"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type Options struct {
	database     string
	hostURI      string
	token        string
	ctx          context.Context
	writeTimeout time.Duration
	queryTimeout time.Duration
	logger       domain.Logger
}

type Option func(*Options)

func WithDatabase(database string) Option {
	return func(opt *Options) {
		opt.database = database
	}
}

func WithHostURL(hostURL string) Option {
	return func(opt *Options) {
		opt.hostURI = hostURL
	}
}

func WithToken(token string) Option {
	return func(opt *Options) {
		opt.token = token
	}
}

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

// NewOptions returns an empty Options struct! REQUIRED to fill it with options
func NewOptions(opts ...Option) *Options {
	return &Options{
		database:     "",
		hostURI:      "",
		token:        "",
		logger:       log.NewSlog(),
		writeTimeout: 10 * time.Second,
		queryTimeout: time.Minute,
		ctx:          context.Background(),
	}
}
