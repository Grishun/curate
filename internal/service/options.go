package service

import (
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/Grishun/curate/internal/storage/memory"
)

type Option func(*Options)

type Options struct {
	logger          domain.Logger
	storage         domain.Storage
	providers       []domain.Provider
	pollingInterval time.Duration
	quote           string
	historyLimit    uint32
}

func WithLogger(l domain.Logger) Option {
	return func(options *Options) {
		options.logger = l
	}
}

func WithStorage(s domain.Storage) Option {
	return func(options *Options) {
		options.storage = s
	}
}

func WithProviders(providers ...domain.Provider) Option {
	return func(options *Options) {
		options.providers = providers
	}
}

func WithPollingInterval(interval time.Duration) Option {
	return func(options *Options) {
		options.pollingInterval = interval
	}
}

func WithQuote(quote string) Option {
	return func(options *Options) {
		options.quote = quote
	}
}

func WithHistoryLimit(limit uint32) Option {
	return func(options *Options) {
		options.historyLimit = limit
	}
}

func NewOptions() *Options {
	return &Options{
		logger:          log.NewSlog(),
		storage:         memory.New(),
		providers:       []domain.Provider{coindesk.New()},
		pollingInterval: time.Minute,
		quote:           "USD",
	}
}
