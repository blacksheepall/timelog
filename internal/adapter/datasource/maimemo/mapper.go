package maimemo

import (
	"math"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

// Metric names produced from MaiMemo's today-items list.
const (
	MetricTodayWords    = "今日学习单词"
	MetricTodayNew      = "今日新学单词"
	MetricTodayFinished = "今日已完成单词"
)

// Metric names produced from MaiMemo's study progress.
const (
	MetricTodayStudyTime = "今日学习时长"
	MetricTodayPlan      = "今日应背单词"
	MetricTodayProgress  = "今日完成率"
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
		{MetricName: MetricTodayWords, Value: float64(len(resp.TodayItems)), Unit: "个", RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayNew, Value: newCount, Unit: "个", RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayFinished, Value: finishedCount, Unit: "个", RecordedAt: now, Source: "maimemo"},
	}
}

// MapStudyProgress converts a MaiMemo study-progress response into domain
// data points: study time in minutes, today's planned word count, and the
// completion percentage.
func MapStudyProgress(resp *GetStudyProgressResponse, now time.Time) []domain.MetricDataPoint {
	if resp == nil {
		return nil
	}

	minutes := round2(float64(resp.Progress.StudyTime) / 60000.0)
	completion := 0.0
	if resp.Progress.Total > 0 {
		completion = round2(float64(resp.Progress.Finished) / float64(resp.Progress.Total) * 100)
	}

	return []domain.MetricDataPoint{
		{MetricName: MetricTodayStudyTime, Value: minutes, Unit: "分钟", RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayPlan, Value: float64(resp.Progress.Total), Unit: "个", RecordedAt: now, Source: "maimemo"},
		{MetricName: MetricTodayProgress, Value: completion, Unit: "%", RecordedAt: now, Source: "maimemo"},
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
