package ai

import (
	"net/http"
	"time"
)

type (
	// Client
	Client struct {
		hc    *http.Client
		token string
	}

	clientOpt func(*Client)
)

const (
	defaultTimeout = 30 * time.Second
)

// NewClient
func NewClient(token string, opts ...clientOpt) *Client {
	c := &Client{
		hc: &http.Client{
			Timeout: defaultTimeout,
		},
		token: token,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithCustomTimeout set custom timeout
func WithCustomTimeout(t time.Duration) clientOpt {
	return func(c *Client) {
		c.hc.Timeout = t
	}
}
