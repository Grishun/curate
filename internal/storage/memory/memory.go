package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/Grishun/curate/internal/domain"
)

type Memory struct {
	data         map[string][]domain.Rate //TODO: change it with linked list
	mu           sync.RWMutex
	historyLimit uint
}

func New() *Memory {
	return &Memory{
		mu:           sync.RWMutex{},
		data:         make(map[string][]domain.Rate),
		historyLimit: 10,
	}
}

func (m *Memory) Get(_ context.Context, currency string, limit uint) ([]domain.Rate, error) {
	if limit > m.historyLimit {
		return nil, errors.New("limit is > than configuired history limit")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	rates := m.data[currency]
	var ratesCopy []domain.Rate // return a copy of rates slice to avoid data races

	if len(rates) > int(limit) {
		ratesCopy = make([]domain.Rate, limit)
		copy(ratesCopy, rates[len(rates)-int(limit):])
	} else {
		ratesCopy = make([]domain.Rate, len(rates))
		copy(ratesCopy, rates)
	}

	return ratesCopy, nil
}

func (m *Memory) GetAll(_ context.Context, limit uint) (map[string][]domain.Rate, error) {
	if limit > m.historyLimit {
		return nil, errors.New("limit is > than configuired history limit")
	}

	// create new map to avoid race condition
	res := make(map[string][]domain.Rate, len(m.data))

	m.mu.RLock()
	for currency, rates := range m.data {
		var ratesCopy []domain.Rate

		if len(rates) > int(limit) {
			ratesCopy = make([]domain.Rate, limit)
			copy(ratesCopy, rates[len(rates)-int(limit):])
		} else {
			ratesCopy = make([]domain.Rate, len(rates))
			copy(ratesCopy, rates)
		}

		res[currency] = ratesCopy
	}
	m.mu.RUnlock()

	return res, nil
}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, newRate := range rates {
		rate, ok := m.data[newRate.Currency]
		if !ok {
			rate = make([]domain.Rate, 0, m.historyLimit)
			m.data[newRate.Currency] = rate
		}
		if len(rate) >= int(m.historyLimit) {
			return errors.New("history limit reached")
		}

		rate = append(rate, newRate)
		m.data[newRate.Currency] = rate
	}

	return nil
}
