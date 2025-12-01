package ratebuffer

import (
	"testing"

	"github.com/Grishun/curate/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestRateBuffer(t *testing.T) {
	lenght := 10

	buf := New(uint(lenght))

	for i := 0; i < lenght; i++ {
		buf.Push(domain.Rate{Value: float64(i)})
	}

	require.EqualValues(t, lenght, buf.Len())

	// validate the last value
	lastRate := buf.LastRate()
	require.EqualValues(t, lastRate.Value, lenght-1)

	// validate all values
	rates := buf.LastNRates(uint(lenght))
	require.Len(t, rates, lenght)

	for i, rate := range rates {
		require.EqualValues(t, rate.Value, i)
	}

	// validate last n values
	n := 5
	LastNRates := buf.LastNRates(uint(n))
	for i := n; i < lenght; i++ {
		require.EqualValues(t, LastNRates[i-n].Value, i)
	}

	// validate rewrite
	for i := 0; i < n; i++ {
		buf.Push(domain.Rate{Value: float64(i)})
	}

	rewrittenRateValues := make([]float64, lenght)
	for _, rate := range buf.LastNRates(uint(lenght)) {
		rewrittenRateValues[int(rate.Value)] = rate.Value
	}

	require.Equal(t, rewrittenRateValues[:5], []float64{0, 1, 2, 3, 4})
	require.Equal(t, rewrittenRateValues[5:], []float64{5, 6, 7, 8, 9})

}
