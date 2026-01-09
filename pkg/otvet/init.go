package otvet

import (
	"net/http"
)

// OtvetClient is the client for otvet.mail.ru API
type OtvetClient struct {
	config     *OtvetClientConfig
	httpClient *http.Client
	baseURL    string
}

// OtvetClientOption is a function type for client options
type OtvetClientOption func(*OtvetClient)

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) OtvetClientOption {
	return func(c *OtvetClient) {
		if baseURL != "" {
			c.baseURL = baseURL
		}
	}
}

// NewOtvetClient creates a new OtvetClient instance
func NewOtvetClient(config *OtvetClientConfig, opts ...OtvetClientOption) (*OtvetClient, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	client := &OtvetClient{
		config:     config,
		httpClient: httpClient,
		baseURL:    OtvetBaseURL,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}
