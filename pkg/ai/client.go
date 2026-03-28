package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL    = "http://127.0.0.1:7111"
	defaultTimeout    = 15 * time.Second
	checkAllergyRoute = "/check-allergy"
)

// Client is an adapter for the external allergy model service.

// BaseURL is loaded from ALLERGEN_FILTER_AI_API_URL when available.
// If the env is empty, it falls back to http://127.0.0.1:7111 for local air runs.
// In Docker Compose, set ALLERGEN_FILTER_AI_API_URL to http://model:7111.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.httpClient.Timeout = timeout
		}
	}
}

func NewClient(baseURL string, opts ...Option) *Client {
	url := strings.TrimSpace(baseURL)
	if url == "" {
		url = defaultBaseURL
	}

	client := &Client{
		baseURL: strings.TrimRight(url, "/"),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func NewClientFromEnv(opts ...Option) *Client {
	return NewClient(os.Getenv("ALLERGEN_FILTER_AI_API_URL"), opts...)
}

// CheckAllergy sends the menu/allergy payload to POST /check-allergy.
// It returns the raw response body for flexible handling in usecases later.
func (c *Client) CheckAllergy(ctx context.Context, payload CheckAllergyRequest) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal check-allergy payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+checkAllergyRoute, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create check-allergy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call check-allergy endpoint: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read check-allergy response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("check-allergy returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
