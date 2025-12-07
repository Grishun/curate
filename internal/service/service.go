package service

import (
	"context"
	"runtime"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/go-co-op/gocron/v2"
	syncmap "github.com/zolstein/sync-map"
)

type Subscription struct {
	Currency string
	Provider *string
}

type Service struct {
	options       *Options
	scheduler     gocron.Scheduler
	subscriptions syncmap.Map[Subscription, chan domain.Rate]
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
		options:       options,
		scheduler:     scheduler,
		subscriptions: syncmap.Map[Subscription, chan domain.Rate]{},
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.options.logger.Debug("starting service")

	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.options.pollingInterval),
		gocron.NewTask(s.fetchAndStore),
		gocron.WithContext(ctx),
	)
	if err != nil {
		s.options.logger.Error("failed to schedule fetch job", "error", err)
		return err
	}

	s.scheduler.Start()
	s.options.logger.Info("fetch scheduler started")

	<-ctx.Done()

	if err := s.Stop(ctx); err != nil {
		s.options.logger.Error("failed to stop service", "error", err)
		return err
	}

	return nil
}

func (s *Service) Stop(_ context.Context) error {
	s.options.logger.Debug("stopping service")

	if err := s.scheduler.Shutdown(); err != nil {
		s.options.logger.Error("failed to shutdown scheduler", "error", err)
		return err
	}

	s.options.logger.Info("scheduler stopped")

	return nil
}

func (s *Service) GetRate(ctx context.Context, currency string, limit uint32) ([]domain.Rate, error) {
	rates, err := s.options.storage.Get(ctx, currency, limit)
	if err != nil {
		s.options.logger.Error("failed to get rate", "currency", currency, "limit", limit, "error", err)
		return nil, err
	}

	s.options.logger.Info("fetched rates for currency", "currency", currency, "count", len(rates))

	return rates, nil
}

func (s *Service) GetRates(ctx context.Context, limit uint32) (map[string][]domain.Rate, error) {
	rates, err := s.options.storage.GetAll(ctx, limit)
	if err != nil {
		s.options.logger.Error("failed to get rates", "limit", limit, "error", err)
		return nil, err
	}

	s.options.logger.Info("fetched rates for all currencies", "limit", limit, "currency_count", len(rates))

	return rates, nil
}

func (s *Service) fetchAndStore(ctx context.Context) {
	for _, provider := range s.options.providers {
		s.options.logger.Debug("fetching rates from provider", "provider", provider.Name())
		result, err := provider.Fetch(ctx)
		if err != nil {
			s.options.logger.Error("failed to fetch provider", "provider", provider.Name(), "error", err)
		}

		s.options.logger.Info("provider fetch finished", "provider", provider.Name(), "rates_count", len(result))

		rates := make([]domain.Rate, 0, len(result))
		for currency, value := range result {
			rate := domain.Rate{
				Currency:  currency,
				Quote:     s.options.quote,
				Provider:  provider.Name(),
				Timestamp: time.Now(),
				Value:     value,
			}

			rates = append(rates, rate)

			s.subscriptions.Range(func(sub Subscription, ch chan domain.Rate) bool {
				if sub.Provider != nil && *sub.Provider != provider.Name() { // if we configuired a provider, we should return only this provider's rates
					return true
				}

				if sub.Currency == currency {
					select {
					case <-ctx.Done():
						return false
					case ch <- rate:
						s.options.logger.Debug("sending rate to subscription", "currency", currency)
					default:
					}
				}

				return true
			})
		}

		s.options.logger.Debug("inserting provider rates", "provider", provider.Name(), "rates", rates)
		err = s.options.storage.Insert(ctx, rates...)
		if err != nil {
			s.options.logger.Error("failed to insert rates to storage", "provider", provider.Name(), "error", err)
		}
	}
}

func (s *Service) SubscribeRate(ctx context.Context, sub Subscription) <-chan domain.Rate {
	subCh, ok := s.subscriptions.LoadOrStore(sub, make(chan domain.Rate))

	if !ok {
		s.options.logger.Debug("new subscription", "currency", sub.Currency)
		go func() {
			<-ctx.Done()
			s.subscriptions.Delete(sub)
			close(subCh)
		}()
	}

	return subCh
}

func (s *Service) HealthCheck(ctx context.Context) (err error) {
	if err := s.options.storage.HealthCheck(ctx); err != nil {
		s.options.logger.Error("service healthcheck failed", "error", err)
		return err
	}

	s.options.logger.Info("service healthcheck passed")

	return nil
}
