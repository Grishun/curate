package service

import (
	"context"
	"runtime"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/go-co-op/gocron/v2"
)

type Service struct {
	options   *Options
	scheduler gocron.Scheduler
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
		options:   options,
		scheduler: scheduler,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.options.logger.Debug("starting service")

	s.scheduler.NewJob(
		gocron.DurationJob(s.options.pollingInterval),
		gocron.NewTask(s.fetchAndStore),
		gocron.WithContext(ctx),
	)

	s.scheduler.Start()
	s.options.logger.Debug("fetch scheduler started")

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Service) Stop(_ context.Context) error {
	s.options.logger.Debug("stopping service")

	return s.scheduler.Shutdown()
}

func (s *Service) GetRate(ctx context.Context, currency string, limit uint) ([]domain.Rate, error) {
	return s.options.storage.Get(ctx, currency)
}

func (s *Service) GetRates(ctx context.Context, limit uint) (map[string][]domain.Rate, error) {
	return s.options.storage.GetAll(ctx)
}

func (s *Service) fetchAndStore(ctx context.Context) {
	for _, provider := range s.options.providers {
		result, err := provider.Fetch(ctx)
		if err != nil {
			s.options.logger.Error("failed to fetch provider", "provider", provider, "error", err)
		}

		rates := make([]domain.Rate, 0, len(result))
		for currency, value := range result {
			rates = append(rates, domain.Rate{
				Currency:  currency,
				Quote:     s.options.quote,
				Provider:  provider.Name(),
				Timestamp: time.Now(),
				Value:     value,
			})
		}

		s.options.logger.Debug("fetched rates from provider. Calling for storage.Insert()", "provider", provider, "rates", rates)
		err = s.options.storage.Insert(ctx, rates...)
		if err != nil {
			s.options.logger.Error("failed to insert rates to storage", "provider", provider, "error", err)
		}
	}
}
