package httplib

import (
	"github.com/go-resty/resty/v2"
	"time"
)

const (
	DefaultMaxRetries   = 3
	DefaultInitialDelay = time.Duration(5) * time.Second
	DefaultMaxDelay     = time.Duration(60) * time.Second
)

type Client struct {
	client *resty.Client
}

type Option func(*Client)

type RetryConfig struct {
	RetryCnt     int
	InitialDelay time.Duration
	MaxWaitTime  time.Duration
}

func NewClient(options ...Option) *resty.Client {
	client := resty.New()

	c := &Client{
		client: client,
	}

	// apply default config first
	c.client.SetRetryCount(DefaultMaxRetries)
	c.client.SetRetryWaitTime(DefaultInitialDelay)
	c.client.SetRetryMaxWaitTime(DefaultMaxDelay)
	c.client.AddRetryCondition(func(response *resty.Response, err error) bool {
		return response.StatusCode() >= 500
	})

	// apply caller passed options
	for _, option := range options {
		option(c)
	}

	return c.client
}

func WithRetryConfig(cfg RetryConfig) Option {
	return func(c *Client) {
		if cfg.RetryCnt == 0 {
			c.client.SetRetryCount(DefaultMaxRetries)
		}

		if cfg.InitialDelay == 0 {
			c.client.SetRetryWaitTime(DefaultInitialDelay)
		}

		if cfg.MaxWaitTime == 0 {
			c.client.SetRetryMaxWaitTime(DefaultMaxDelay)
		}
	}
}

func WithRetryCondition(condition resty.RetryConditionFunc) Option {
	return func(c *Client) {
		c.client.AddRetryCondition(condition)
	}
}
