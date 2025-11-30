package memory

import (
	"context"
	"sync"

	"github.com/Grishun/curate/internal/domain"
)

type Memory struct {
	data map[string][]domain.Rate //TODO: change it with linked list
	mu   sync.RWMutex
}

func New() *Memory {
	return &Memory{
		mu:   sync.RWMutex{},
		data: make(map[string][]domain.Rate),
	}
}

func (m *Memory) Get(_ context.Context, currency string) ([]domain.Rate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data[currency], nil
}

func (m *Memory) GetAll(_ context.Context) (map[string][]domain.Rate, error) {
	// create new map
	res := make(map[string][]domain.Rate, len(m.data))

	for currency, rates := range m.data {
		res[currency] = rates
	}

	return res, nil
}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, newRate := range rates {
		rate, ok := m.data[newRate.Currency]
		if !ok {
			rate = make([]domain.Rate, 0)
			m.data[newRate.Currency] = rate
		}

		rate = append(rate, newRate)
		m.data[newRate.Currency] = rate
	}

	return nil
}
