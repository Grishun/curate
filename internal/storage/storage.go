package storage

import (
	"context"
	"github.com/Grishun/curate/internal/domain"
)

type Storage interface {
	Get(ctx context.Context, code string) (domain.Rate, error)
	GetAll(ctx context.Context) ([]domain.Rate, error)

	Insert(ctx context.Context, rates ...domain.Rate) error
}
