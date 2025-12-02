package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/Grishun/curate/internal/domain"
)

func NewOptions() *domain.RequestOptions {
	return &domain.RequestOptions{
		Headers:     make(http.Header),
		QueryParams: make(url.Values),
		Body:        bytes.NewBuffer(nil),
		URI:         "",
		Method:      "",
	}
}

func WithHeaders(headers http.Header) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.Headers = headers
	}
}

func WithQueryParams(queryParams url.Values) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.QueryParams = queryParams
	}
}

func WithBody(body io.Reader) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.Body = body
	}
}

func WithURI(uri string) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.URI = uri
	}
}

func WithMethod(method string) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.Method = method
	}
}

// WithUnmarshallTo is used to unmarshal the response body to a specific type. Transfer ONLY pointer types
func WithUnmarshallTo(unmarshallTo any) domain.RequestOption {
	return func(opt *domain.RequestOptions) {
		opt.UnmarshallTo = unmarshallTo
	}
}
