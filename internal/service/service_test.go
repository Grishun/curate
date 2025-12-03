package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
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
		storageMap, err := service.options.storage.GetAll(ctx, 10)

		return err == nil && len(storageMap) == 3
	}, time.Second*3, time.Second)

	cancel()
	require.NoError(t, <-errCh)
}
