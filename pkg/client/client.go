package client

import (
	"context"
	"net/http"
)

// Client defines the HTTP client interface used by Providers.
type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// HTTPClient implements Client using standard net/http and custom RoundTripper middleware.
type HTTPClient struct {
	HTTPClient *http.Client
}

func NewHTTPClient(transport http.RoundTripper) *HTTPClient {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &HTTPClient{
		HTTPClient: &http.Client{Transport: transport},
	}
}

func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	return c.HTTPClient.Do(req)
}
