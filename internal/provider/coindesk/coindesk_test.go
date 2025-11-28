package coindesk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoindeskProvider(t *testing.T) {
	mockedCoinDesk := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"BTC":{"USD":91194.7},"ETH":{"USD":3050.22},"TRX":{"USD":0.2813}}`))
		w.WriteHeader(http.StatusOK)
	}))

	defer mockedCoinDesk.Close()

	provider := New(WithURI(mockedCoinDesk.URL))
	require.Equal(t, "127.0.0.1", provider.Name())

	responseMap, err := provider.Fetch(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, responseMap)

	for currency, rate := range responseMap {
		t.Logf("currency %s rate %f", currency, rate)
		require.NotZero(t, rate)
		require.NotZero(t, currency)
	}
}
