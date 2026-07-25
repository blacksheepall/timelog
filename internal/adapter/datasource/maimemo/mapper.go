package maimemo

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

// MapTodayItems converts a MaiMemo response into domain data points.
// The source identifier is hard-coded to the datasource name "maimemo".
func MapTodayItems(resp *GetTodayItemsResponse, now time.Time) []domain.MetricDataPoint {
	if resp == nil || len(resp.Items) == 0 {
		return nil
	}

	points := make([]domain.MetricDataPoint, 0, len(resp.Items))
	for _, item := range resp.Items {
		name := item.Name
		if name == "" {
			continue
		}
		points = append(points, domain.MetricDataPoint{
			MetricName: name,
			Value:      item.Value,
			RecordedAt: now,
			Source:     "maimemo",
		})
	}
	return points
}
