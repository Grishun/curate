package domain

import "time"

type HistoryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type Rate struct {
	Currency string `json:"code"`
	Quote    string `json:"quote"`
	Provider string `json:"provider"`
	//History   []HistoryPoint `json:"history"` // TODO: change it with linked list
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	//Metadata  map[string]any `json:"metadata"`
}
