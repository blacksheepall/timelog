package maimemo

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

// Metric names produced from MaiMemo's today-items list.
const (
	MetricTodayWords    = "今日学习单词"
	MetricTodayNew      = "今日新学单词"
	MetricTodayFinished = "今日已完成单词"
)

// MapTodayItems converts a MaiMemo response into domain data points.
// The API returns today's word list, so the metrics are aggregate counts.
// An empty list means either "studied nothing" or "App not opened today"
// (the API cannot distinguish the two), so no points are emitted then.
// The source identifier is hard-coded to the datasource name "maimemo".
func MapTodayItems(resp *GetTodayItemsResponse, now time.Time) []domain.MetricDataPoint {
	if resp == nil || len(resp.TodayItems) == 0 {
		return nil
	}

	var newCount, finishedCount float64
	for _, item := range resp.TodayItems {
		if item.IsNew {
			newCount++
		}
		if item.IsFinished {
			finishedCount++
		}
	}

	return []domain.MetricDataPoint{
		{MetricName: MetricTodayWords, Value: float64(len(resp.TodayItems)), RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayNew, Value: newCount, RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayFinished, Value: finishedCount, RecordedAt: now, Source: "maimemo"},
	}
}
