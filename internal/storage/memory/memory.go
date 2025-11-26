package memory

import (
	"context"
	"github.com/Grishun/curate/internal/domain"
	"sync"
)

type Memory struct {
	data map[string]domain.Rate
	opts *Options
	mu   sync.RWMutex
}

var mockedData = map[string]domain.Rate{
	"BTC": {
		Code:     "BTC",
		Quote:    "USD",
		Provider: "coindesk",
		History: []domain.HistoryPoint{
			{Value: 100.1},
			{Value: 100.2},
			{Value: 100.3},
		},
	},
	"ETH": {
		Code:     "ETH",
		Quote:    "USD",
		Provider: "coindesk",
		History: []domain.HistoryPoint{
			{Value: 400.1},
			{Value: 400.2},
			{Value: 400.3},
		},
	},
}

func New(opts ...Option) *Memory {
	m := Memory{
		mu:   sync.RWMutex{},
		opts: NewOptions(),
		data: mockedData,
	}

	for _, opt := range opts {
		opt(m.opts)
	}

	return &m
}

func (m *Memory) Get(_ context.Context, code string) (domain.Rate, error) {
	return m.data[code], nil
}

func (m *Memory) GetAll(_ context.Context) ([]domain.Rate, error) {
	res := make([]domain.Rate, 0, len(m.data))

	for _, rate := range m.data {
		res = append(res, rate)
	}

	return res, nil
}

func (m *Memory) Insert(_ context.Context, rates ...domain.Rate) error {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	//
	//next := make(map[string]domain.Rate, len(rates))
	//
	//for _, rate := range rates {
	//	previous, ok := m.data[rate.Code]
	//	if !ok {
	//		previous = domain.Rate{}
	//	}
	//
	//	historyPoint := domain.HistoryPoint{
	//		Timestamp: rate.UpdatedAt,
	//		Value:     rate.Value,
	//	}
	//
	//	history := append([]domain.HistoryPoint{historyPoint}, previous.History...)
	//
	//	rate.History = history
	//}
	//
	//m.data = next

	return nil
}
