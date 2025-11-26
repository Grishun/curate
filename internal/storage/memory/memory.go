package memory

import (
	"context"
	"github.com/Grishun/curate/internal/domain"
	"sync"
)

type Memory struct {
	data map[string]domain.Rate
	mu   sync.RWMutex
}

func New() *Memory {
	return &Memory{
		data: make(map[string]domain.Rate),
		mu:   sync.RWMutex{},
	}
}

func (m *Memory) Get(_ context.Context, code string) (domain.Rate, error) {

}

func (m *Memory) GetAll(_ context.Context) ([]domain.Rate, error) {

}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	next := make(map[string]domain.Rate, len(rates))

	for _, rate := range rates {
		previous, ok := m.data[rate.Code]
		if !ok {
			previous = domain.Rate{}
		}

		historyPoint := domain.HistoryPoint{
			Timestamp: rate.UpdatedAt,
			Value:     rate.Value,
		}

		history := append([]domain.HistoryPoint{historyPoint}, previous.History...)

		rate.History = history
	}

	m.data = next

	return nil
}
