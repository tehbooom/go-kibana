package kbapi

import (
	"net/http"
)

// Interface for performing requests (allows mocking for tests)
type Transport interface {
	Perform(req *http.Request) (*http.Response, error)
}

// Option defines a functional option for configuring API requests
type RequestOption func(*http.Request) error

// WithHeaders adds custom headers to the request
func WithHeaders(headers http.Header) RequestOption {
	return func(req *http.Request) error {
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		return nil
	}
}
