package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Grishun/curate/internal/domain"
)

type Memory struct {
	data map[string]*domain.Rate //TODO: change it with linked list
	mu   sync.RWMutex
}

func NewMemoryStorage() *Memory {
	return &Memory{
		mu:   sync.RWMutex{},
		data: make(map[string]*domain.Rate),
	}
}

func (m *Memory) Get(_ context.Context, currency string) (domain.Rate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return *m.data[currency], nil
}

func (m *Memory) GetAll(_ context.Context) ([]domain.Rate, error) {
	res := make([]domain.Rate, 0, len(m.data))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rates := range m.data {
		res = append(res, *rates)
	}

	return res, nil
}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, newRate := range rates {
		rate, ok := m.data[newRate.Currency]
		if !ok {
			rate = &domain.Rate{}
			m.data[newRate.Currency] = rate
		}

		rate.Currency = newRate.Currency
		rate.Quote = newRate.Quote
		rate.Provider = newRate.Provider
		rate.Timestamp = newRate.Timestamp
		rate.Value = newRate.Value
		rate.History = append(rate.History, domain.HistoryPoint{
			Timestamp: time.Now(),
			Value:     newRate.Value,
		})
		newRate.Value = rate.Value
	}

	return nil
}
