package memory

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorageGet(t *testing.T) {
	storage := New()

	currencies := []string{"BTC", "ETH", "LUNA"}
	ratesQty := 10

	for _, currency := range currencies {
		storage.Insert(context.Background(), generateTestRates(currency, ratesQty)...)
	}

	limit := 5
	for _, currency := range currencies {
		rates, err := storage.Get(context.Background(), currency, uint(limit))
		require.NoError(t, err)

		require.Len(t, rates, limit)
		for _, rate := range rates {
			require.NotZero(t, rate.Value)
			require.NotZero(t, rate.Timestamp)
			require.Equal(t, rate.Provider, "https://min-api.cryptocompare.com")
			require.Equal(t, rate.Quote, "USD")
		}
	}
}

func TestMemoryStorageGetAll(t *testing.T) {
	storage := New()

	currencies := []string{"BTC", "ETH", "LUNA"}
	ratesQty := 10

	for _, currency := range currencies {
		err := storage.Insert(context.Background(), generateTestRates(currency, ratesQty)...)
		require.NoError(t, err)
	}

	limit := 5
	ratesMap, err := storage.GetAll(context.Background(), uint(limit))
	require.NoError(t, err)
	require.Len(t, ratesMap, len(currencies))

	for currency, rates := range ratesMap {
		require.Contains(t, currencies, currency)
		require.Len(t, rates, limit)

		for _, rate := range rates {
			require.NotZero(t, rate.Value)
			require.NotZero(t, rate.Timestamp)
			require.Equal(t, rate.Provider, "https://min-api.cryptocompare.com")
			require.Equal(t, rate.Quote, "USD")
		}
	}
}

func generateTestRates(currency string, qty int) []domain.Rate {
	res := make([]domain.Rate, qty)

	for i := 0; i < qty; i++ {
		res[i] = domain.Rate{
			Currency:  currency,
			Quote:     "USD",
			Provider:  "https://min-api.cryptocompare.com",
			Value:     rand.Float64(),
			Timestamp: time.Now(),
		}
	}

	return res
}
