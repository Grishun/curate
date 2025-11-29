package service

import (
	"runtime"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/Grishun/curate/internal/storage/memory"
	"github.com/go-co-op/gocron/v2"
)

type Option func(*Options)

type Options struct {
	logger          domain.Logger
	storage         domain.Storage
	providers       []domain.Provider
	pollingInterval time.Duration
	quote           string
	scheduler       gocron.Scheduler
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

func WithScheduler(scheduler gocron.Scheduler) Option {
	return func(options *Options) {
		options.scheduler = scheduler
	}
}

func NewOptions() *Options {
	logger := log.NewSlog()

	scheduler, err := gocron.NewScheduler(
		gocron.WithLogger(logger),
		gocron.WithLimitConcurrentJobs(uint(runtime.GOMAXPROCS(0)), gocron.LimitModeReschedule),
	)
	if err != nil {
		logger.Error("failed to create scheduler", "error", err)
	}

	return &Options{
		logger:          log.NewSlog(),
		storage:         memory.NewMemoryStorage(),
		providers:       []domain.Provider{coindesk.New()},
		pollingInterval: time.Minute,
		quote:           "USD",
		scheduler:       scheduler,
	}
}
