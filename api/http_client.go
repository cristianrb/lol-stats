package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(endpoint string, response any) error
}

type DefaultHTTPClient struct {
	apiKey string
	client *http.Client
}

const ApiKeyHeader = "X-Riot-Token"

func NewHTTPClient(apiKey string) *DefaultHTTPClient {
	return &DefaultHTTPClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *DefaultHTTPClient) Do(endpoint string, response any) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set(ApiKeyHeader, c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &response)
}
