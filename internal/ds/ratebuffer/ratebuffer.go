package ratebuffer

import (
	"container/ring"

	"github.com/Grishun/curate/internal/domain"
)

// TODO: add comments

type Buffer interface {
	Push(rate domain.Rate)
	LastNRates(n uint) []domain.Rate
	LastRate() domain.Rate
	Len() uint
}

type RateBuffer struct {
	ring *ring.Ring
}

func New(len uint) *RateBuffer {
	return &RateBuffer{
		ring: ring.New(int(len)),
	}
}

func (rb *RateBuffer) Push(rate domain.Rate) {
	rb.ring.Value = rate
	rb.ring = rb.ring.Move(1)
}

func (rb *RateBuffer) LastNRates(n uint) []domain.Rate {
	if rb.Len() < n {
		n = rb.Len()
	}

	result := make([]domain.Rate, n)

	rb.ring = rb.ring.Move(-int(n))

	for i := 0; i < int(n); i++ {
		result[i] = rb.ring.Value.(domain.Rate)
		rb.ring = rb.ring.Next()
	}

	return result
}

func (rb *RateBuffer) LastRate() domain.Rate {
	return rb.ring.Prev().Value.(domain.Rate)
}

func (rb *RateBuffer) Len() uint {
	return uint(rb.ring.Len())
}
