package memory

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	currencies := []string{"BTC", "ETH", "LUNA"}
	ratesQty := 5

	for _, currency := range currencies {
		storage.Insert(context.Background(), generateTestRates(currency, ratesQty)...)
	}

	for _, currency := range currencies {
		rate, err := storage.Get(context.Background(), currency)
		require.NoError(t, err)

		require.NotZero(t, rate.Value)
		require.NotZero(t, rate.Timestamp)
		require.Equal(t, currency, rate.Currency)
		require.Equal(t, "USD", rate.Quote)
		require.Equal(t, "https://min-api.cryptocompare.com", rate.Provider)

		require.Len(t, rate.History, ratesQty)
		for _, historyPoint := range rate.History {
			require.NotZero(t, historyPoint.Value)
			require.NotZero(t, historyPoint.Timestamp)
		}

		require.Equal(t, rate.History[ratesQty-1].Value, rate.Value)

		require.Len(t, storage.data[currency].History, ratesQty)
	}
}

func generateTestRates(currency string, qty int) []domain.Rate {
	res := make([]domain.Rate, qty)

	for i := 0; i < qty; i++ {
		res[i] = domain.Rate{
			Currency:  currency,
			Quote:     "USD",
			Provider:  "https://min-api.cryptocompare.com",
			History:   []domain.HistoryPoint{},
			Value:     rand.Float64(),
			Timestamp: time.Now(),
		}
	}

	return res
}
