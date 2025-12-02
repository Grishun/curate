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
	Currency  string    `lp:"measuremnt" json:"currency"`
	Quote     string    `lp:"tag,quote" json:"quote"`
	Provider  string    `lp:"tag,provider" json:"provider"`
	Value     float64   `lp:"field,value" json:"value"`
	Timestamp time.Time `lp:"timestamp" json:"timestamp"`
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
