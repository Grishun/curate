package domain

import (
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
	Currency  string    `lp:"measurement" json:"currency"`
	Quote     string    `lp:"tag,quote" json:"quote"`
	Provider  string    `lp:"tag,provider" json:"provider"`
	Value     float64   `lp:"field,value" json:"value"` // TODO: use decimal instead
	Timestamp time.Time `lp:"timestamp" json:"timestamp"`
	//Metadata  map[string]any `json:"metadata"`
}
type RequestOption func(opt *RequestOptions)

type RequestOptions struct {
	Headers      http.Header
	QueryParams  url.Values
	Body         io.Reader
	URI          string
	Method       string
	UnmarshallTo any
}
