package domain

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

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

type RequestOption func(opt *RequestOptions)

type RequestOptions struct {
	Ctx          context.Context
	Headers      http.Header
	QueryParams  url.Values
	Body         io.Reader
	URI          string
	Method       string
	UnmarshallTo any
}
