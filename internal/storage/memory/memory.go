package memory

import (
	"context"
	"sync"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/ds/ratebuffer"
)

type Memory struct {
	data         map[string]ratebuffer.Buffer
	mu           sync.RWMutex
	historyLimit uint
}

func New() *Memory {
	return &Memory{
		mu:           sync.RWMutex{},
		data:         make(map[string]ratebuffer.Buffer),
		historyLimit: 10,
	}
}

func (m *Memory) GetHistoryLimit() uint {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.historyLimit
}

func (m *Memory) Get(_ context.Context, currency string, limit uint) ([]domain.Rate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data[currency].LastNRates(limit), nil
}

func (m *Memory) GetAll(_ context.Context, limit uint) (map[string][]domain.Rate, error) {
	// create a new map to avoid race condition
	res := make(map[string][]domain.Rate, len(m.data))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for currency, rates := range m.data {
		res[currency] = rates.LastNRates(limit)
	}

	return res, nil
}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, newRate := range rates {
		rateBuffer, ok := m.data[newRate.Currency]
		if !ok {
			rateBuffer = ratebuffer.New(m.historyLimit)
			m.data[newRate.Currency] = rateBuffer
		}

		rateBuffer.Push(newRate)
		m.data[newRate.Currency] = rateBuffer
	}

	return nil
}
