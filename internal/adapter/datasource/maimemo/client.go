package maimemo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// API paths from the official MaiMemo Open API bundle
// (https://open.maimemo.com/api_bundle.yaml). The endpoint host is
// configurable; the API itself lives under the /open prefix on that host.
const (
	getTodayItemsPath    = "/open/api/v1/memo/study/get_today_items"
	getStudyProgressPath = "/open/api/v1/memo/study/get_study_progress"
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

// GetTodayItemsResponse mirrors the data payload of MaiMemo's GetTodayItems.
type GetTodayItemsResponse struct {
	TodayItems []TodayItem `json:"today_items"`
}

// TodayItem represents one word in today's study list.
type TodayItem struct {
	VocID       string `json:"voc_id"`
	VocSpelling string `json:"voc_spelling"`
	Order       int    `json:"order"`
	IsNew       bool   `json:"is_new"`
	IsFinished  bool   `json:"is_finished"`
}

// GetStudyProgressResponse mirrors the data payload of MaiMemo's GetStudyProgress.
type GetStudyProgressResponse struct {
	Progress StudyProgress `json:"progress"`
}

// StudyProgress is today's study progress; StudyTime is in milliseconds.
type StudyProgress struct {
	Finished  int   `json:"finished"`
	Total     int   `json:"total"`
	StudyTime int64 `json:"study_time"`
}

// apiEnvelope is the common MaiMemo response wrapper.
type apiEnvelope struct {
	Success bool            `json:"success"`
	Errors  []apiError      `json:"errors"`
	Data    json.RawMessage `json:"data"`
}

type apiError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// GetTodayItems fetches today's study items.
func (c *Client) GetTodayItems(ctx context.Context) (*GetTodayItemsResponse, error) {
	// The default limit is 50; raise it so heavy study days are not truncated.
	var body GetTodayItemsResponse
	if err := c.post(ctx, getTodayItemsPath, []byte(`{"limit":1000}`), &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// GetStudyProgress fetches today's study progress (finished/total/study time).
func (c *Client) GetStudyProgress(ctx context.Context) (*GetStudyProgressResponse, error) {
	var body GetStudyProgressResponse
	if err := c.post(ctx, getStudyProgressPath, []byte(`{}`), &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// post issues a POST with a JSON body and decodes the MaiMemo envelope into out.
func (c *Client) post(ctx context.Context, path string, payload []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("maimemo api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("maimemo api returned status %d", resp.StatusCode)
	}

	var env apiEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return fmt.Errorf("decode maimemo response: %w", err)
	}
	if !env.Success {
		msgs := make([]string, 0, len(env.Errors))
		for _, e := range env.Errors {
			msgs = append(msgs, e.Code+": "+e.Msg)
		}
		return fmt.Errorf("maimemo api error: %s", strings.Join(msgs, "; "))
	}

	if err := json.Unmarshal(env.Data, out); err != nil {
		return fmt.Errorf("decode maimemo data: %w", err)
	}
	return nil
}
