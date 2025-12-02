package http

import (
	"github.com/Grishun/curate/internal/domain"
)

func mapDomainRateToOpenAPIRate(rates []domain.Rate) *Rate {
	if len(rates) == 0 {
		return nil
	}

	lastIndex := len(rates) - 1 // get the last value as it's the newest

	openAPIRate := &Rate{
		Currency:  rates[0].Currency,
		Provider:  rates[0].Provider,
		Quote:     rates[0].Quote,
		Timestamp: rates[lastIndex].Timestamp,
		Value:     rates[lastIndex].Value,
	}

	history := make([]HistoryPoint, len(rates))
	for i, rate := range rates {
		history[i] = HistoryPoint{
			Timestamp: rate.Timestamp,
			Value:     rate.Value,
		}
	}

	openAPIRate.History = history

	return openAPIRate
}

func mapDomainRatesToOpenAPIRates(ratesMap map[string][]domain.Rate) []Rate {
	result := make([]Rate, 0, len(ratesMap))

	for _, rates := range ratesMap {
		result = append(result, *mapDomainRateToOpenAPIRate(rates))
	}

	return result
}
