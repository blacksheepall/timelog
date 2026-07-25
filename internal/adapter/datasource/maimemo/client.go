package maimemo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client calls the MaiMemo Open API.
type Client struct {
	endpoint string
	token    string
	client   *http.Client
}

// NewClient creates a MaiMemo API client.
func NewClient(endpoint, token string) *Client {
	if endpoint == "" {
		endpoint = "https://open.maimemo.com"
	}
	return &Client{
		endpoint: endpoint,
		token:    token,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// GetTodayItemsResponse mirrors the shape of MaiMemo's GetTodayItems response.
// Adjust fields once the actual API schema is confirmed.
type GetTodayItemsResponse struct {
	Items []TodayItem `json:"items"`
}

// TodayItem represents one countable item returned by GetTodayItems.
type TodayItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// GetTodayItems fetches today's study items.
func (c *Client) GetTodayItems(ctx context.Context) (*GetTodayItemsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/maimemo.openapi.study.v1.StudyService/GetTodayItems", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("maimemo api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("maimemo api returned status %d", resp.StatusCode)
	}

	var body GetTodayItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode maimemo response: %w", err)
	}
	return &body, nil
}
