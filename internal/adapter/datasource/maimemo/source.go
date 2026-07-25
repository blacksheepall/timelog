package maimemo

import (
	"context"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/domain"
)

// DataSource implements ports.DataSource for MaiMemo.
type DataSource struct {
	name   string
	client *Client
}

// NewDataSource creates a MaiMemo DataSource from config.
func NewDataSource(name string, cfg config.DatasourceConfig) (*DataSource, error) {
	if name == "" {
		return nil, fmt.Errorf("maimemo datasource name is required")
	}

	rawToken, _ := cfg.Config["token"].(string)
	rawEndpoint, _ := cfg.Config["endpoint"].(string)
	if rawToken == "" {
		return nil, fmt.Errorf("maimemo datasource %q: token is required", name)
	}

	return &DataSource{
		name:   name,
		client: NewClient(rawEndpoint, rawToken),
	}, nil
}

// Name returns the configured datasource name.
func (s *DataSource) Name() string {
	return s.name
}

// Fetch retrieves today's items from MaiMemo and maps them to domain points.
func (s *DataSource) Fetch(ctx context.Context) ([]domain.MetricDataPoint, error) {
	resp, err := s.client.GetTodayItems(ctx)
	if err != nil {
		return nil, err
	}
	return MapTodayItems(resp, time.Now().UTC()), nil
}
