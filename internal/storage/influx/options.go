package influx

import (
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type Options struct {
	writeTimeout time.Duration
	queryTimeout time.Duration
	logger       domain.Logger
	currencies   []string
	database     string
	hostURI      string
	token        string
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

func WithCurrencies(currencies ...string) Option {
	return func(opt *Options) {
		opt.currencies = currencies
	}
}

func WithDatabase(database string) Option {
	return func(opt *Options) {
		opt.database = database
	}
}

func WithHostURI(uri string) Option {
	return func(opt *Options) {
		opt.hostURI = uri
	}
}

func WithToken(token string) Option {
	return func(opt *Options) {
		opt.token = token
	}
}

// NewOptions returns an empty Options struct! REQUIRED to fill it with options
func NewOptions(opts ...Option) *Options {
	return &Options{
		logger:       log.NewSlog(),
		writeTimeout: 10 * time.Second,
		queryTimeout: time.Minute,
		currencies:   []string{"BTC", "ETH", "TRX"},
		database:     "curate",
		hostURI:      "127.0.0.1:8181",
		token:        "",
	}
}
