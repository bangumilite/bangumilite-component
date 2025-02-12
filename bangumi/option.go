package bangumi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func WithAccessToken(token string) RequestOption {
	return func(req *resty.Request) {
		if token != "" {
			req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		}
	}
}

func applyRequestOptions(req *resty.Request, opts ...RequestOption) {
	for _, opt := range opts {
		opt(req)
	}
}
