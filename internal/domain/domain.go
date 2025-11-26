package domain

import "time"

type HistoryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type Rate struct {
	Code     string         `json:"code"`
	Quote    string         `json:"quote"`
	Provider string         `json:"provider"`
	History  []HistoryPoint `json:"history"`
	Metadata map[string]any `json:"metadata"`
}
