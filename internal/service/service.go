package service

import (
	"context"
	"runtime"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/go-co-op/gocron/v2"
)

type Service struct {
	logger          domain.Logger
	storage         domain.Storage
	providers       []domain.Provider
	pollingInterval time.Duration
	quote           string
	scheduler       gocron.Scheduler
}

func New(opts ...Option) *Service {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLogger(options.logger),
		gocron.WithLimitConcurrentJobs(uint(runtime.GOMAXPROCS(0)), gocron.LimitModeReschedule),
	)
	if err != nil {
		options.logger.Error("failed to create scheduler", "error", err)
	}

	return &Service{
		logger:          options.logger,
		storage:         options.storage,
		providers:       options.providers,
		pollingInterval: options.pollingInterval,
		quote:           options.quote,
		scheduler:       scheduler,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.logger.Debug("starting service")

	s.scheduler.NewJob(
		gocron.DurationJob(s.pollingInterval),
		gocron.NewTask(s.fetchAndStore),
		gocron.WithContext(ctx),
	)

	s.scheduler.Start()
	s.logger.Debug("fetch scheduler started")

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Service) Stop(_ context.Context) error {
	s.logger.Debug("stopping service")

	return s.scheduler.Shutdown()
}

func (s *Service) GetRate(ctx context.Context, currency string) ([]domain.Rate, error) {
	return s.storage.Get(ctx, currency)
}

func (s *Service) GetRates(ctx context.Context) (map[string][]domain.Rate, error) {
	return s.storage.GetAll(ctx)
}

func (s *Service) fetchAndStore(ctx context.Context) {
	for _, provider := range s.providers {
		result, err := provider.Fetch(ctx)
		if err != nil {
			s.logger.Error("failed to fetch provider", "provider", provider, "error", err)
		}

		rates := make([]domain.Rate, 0, len(result))
		for currency, value := range result {
			rates = append(rates, domain.Rate{
				Currency:  currency,
				Quote:     s.quote,
				Provider:  provider.Name(),
				Timestamp: time.Now(),
				Value:     value,
			})
		}

		s.logger.Debug("fetched rates from provider. Calling for storage.Insert()", "provider", provider, "rates", rates)
		err = s.storage.Insert(ctx, rates...)
		if err != nil {
			s.logger.Error("failed to insert rates to storage", "provider", provider, "error", err)
		}
	}

}
