package maimemo

import (
	"testing"
	"time"
)

func TestMapTodayItems(t *testing.T) {
	resp := &GetTodayItemsResponse{
		Items: []TodayItem{
			{Name: "今日已背单词", Value: 120, Unit: "个"},
			{Name: "连续打卡天数", Value: 45, Unit: "天"},
		},
	}
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)

	points := MapTodayItems(resp, now)
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if points[0].MetricName != "今日已背单词" || points[0].Value != 120 {
		t.Fatalf("unexpected first point: %+v", points[0])
	}
	if points[1].Source != "maimemo" {
		t.Fatalf("expected source maimemo, got %q", points[1].Source)
	}
	if !points[0].RecordedAt.Equal(now) {
		t.Fatalf("expected recorded_at %v, got %v", now, points[0].RecordedAt)
	}
}

func TestMapTodayItems_SkipsEmptyName(t *testing.T) {
	resp := &GetTodayItemsResponse{
		Items: []TodayItem{
			{Name: "", Value: 1},
			{Name: "有效指标", Value: 2},
		},
	}
	points := MapTodayItems(resp, time.Now())
	if len(points) != 1 || points[0].MetricName != "有效指标" {
		t.Fatalf("expected only valid item, got %+v", points)
	}
}

func TestMapTodayItems_Nil(t *testing.T) {
	points := MapTodayItems(nil, time.Now())
	if points != nil {
		t.Fatalf("expected nil, got %+v", points)
	}
}
