package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/clients/rest"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/stretchr/testify/require"
)

func TestServiceWithMockedCoindesk(t *testing.T) {
	mockedCoinDesk := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"BTC":{"USD":91194.7},"ETH":{"USD":3050.22},"TRX":{"USD":0.2813}}`))
	}))
	defer mockedCoinDesk.Close()

	provider := coindesk.New(coindesk.WithURI(mockedCoinDesk.URL))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := New(WithProviders(provider), WithPollingInterval(time.Second))

	errCh := make(chan error)
	go func() {
		errCh <- service.Start(ctx)
	}()

	require.Eventually(t, func() bool {
		storageMap, err := service.options.storage.GetAll(ctx, service.options.storage.GetHistoryLimit()) // TODO: make more checks
		if err != nil {
			t.Logf("error while getting rates from storage: %v", err)
		}
		return err == nil && len(storageMap) == 3
	}, time.Second*3, time.Second)

	cancel()
	require.NoError(t, <-errCh)
}

func TestServiceWithRealCoindesk(t *testing.T) {
	logger := log.NewSlog()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	provider := coindesk.New(
		coindesk.WithURI("https://min-api.cryptocompare.com"),
		coindesk.WithLogger(logger),
		coindesk.WithHTTPClient(rest.NewClient(rest.WithLogger(logger), rest.WithContext(ctx))),
		coindesk.WithToken("USD"),
		coindesk.WithCurrencies("BTC", "ETH", "TRX"),
	)

	service := New(
		WithProviders(provider),
		WithPollingInterval(time.Second),
	)

	errCh := make(chan error)
	go func() {
		errCh <- service.Start(ctx)
	}()

	t.Run("valid limit", func(t *testing.T) {
		require.Eventually(t, func() bool {
			storageMap, err := service.options.storage.GetAll(ctx, service.options.storage.GetHistoryLimit()) // TODO: make more checks
			if err != nil {
				service.options.logger.Error("error while getting rates from storage", err, err.Error())
			}
			return err == nil && len(storageMap) == 3
		}, time.Second*3, time.Second)
	})

	t.Run("invalid limit", func(t *testing.T) {
		require.Eventually(t, func() bool {
			_, err := service.options.storage.GetAll(ctx, service.options.storage.GetHistoryLimit()+1) // TODO: make more checks
			if err != nil {
				service.options.logger.Error("error while getting rates from storage", err, err.Error())
			}
			return err != nil
		}, time.Second*3, time.Second)
	})

	cancel()
	require.NoError(t, <-errCh)
}
