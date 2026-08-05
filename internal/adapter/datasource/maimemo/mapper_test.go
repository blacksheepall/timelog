package maimemo

import (
	"testing"
	"time"
)

func TestMapTodayItems(t *testing.T) {
	resp := &GetTodayItemsResponse{
		TodayItems: []TodayItem{
			{VocID: "1", VocSpelling: "apple", Order: 1, IsNew: true, IsFinished: true},
			{VocID: "2", VocSpelling: "banana", Order: 2, IsNew: false, IsFinished: true},
			{VocID: "3", VocSpelling: "cherry", Order: 3, IsNew: false, IsFinished: false},
		},
	}
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)

	points := MapTodayItems(resp, now)
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}

	want := map[string]float64{
		MetricTodayWords:    3,
		MetricTodayNew:      1,
		MetricTodayFinished: 2,
	}
	for _, p := range points {
		v, ok := want[p.MetricName]
		if !ok {
			t.Fatalf("unexpected metric name %q", p.MetricName)
		}
		if p.Value != v {
			t.Fatalf("metric %q: expected %v, got %v", p.MetricName, v, p.Value)
		}
		if p.Source != "maimemo" {
			t.Fatalf("expected source maimemo, got %q", p.Source)
		}
		if !p.RecordedAt.Equal(now) {
			t.Fatalf("expected recorded_at %v, got %v", now, p.RecordedAt)
		}
	}
}

func TestMapTodayItems_Empty(t *testing.T) {
	resp := &GetTodayItemsResponse{TodayItems: []TodayItem{}}
	if points := MapTodayItems(resp, time.Now()); points != nil {
		t.Fatalf("expected nil for empty list, got %+v", points)
	}
}

func TestMapTodayItems_Nil(t *testing.T) {
	points := MapTodayItems(nil, time.Now())
	if points != nil {
		t.Fatalf("expected nil, got %+v", points)
	}
}

func TestMapStudyProgress(t *testing.T) {
	resp := &GetStudyProgressResponse{
		Progress: StudyProgress{Finished: 10, Total: 20, StudyTime: 114514},
	}
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)

	points := MapStudyProgress(resp, now)
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}

	want := map[string]float64{
		MetricTodayStudyTime: 1.91, // 114514ms -> 1.91 minutes
		MetricTodayPlan:      20,
		MetricTodayProgress:  50,
	}
	for _, p := range points {
		v, ok := want[p.MetricName]
		if !ok {
			t.Fatalf("unexpected metric name %q", p.MetricName)
		}
		if p.Value != v {
			t.Fatalf("metric %q: expected %v, got %v", p.MetricName, v, p.Value)
		}
		if !p.RecordedAt.Equal(now) {
			t.Fatalf("expected recorded_at %v, got %v", now, p.RecordedAt)
		}
	}
}

func TestMapStudyProgress_ZeroTotal(t *testing.T) {
	resp := &GetStudyProgressResponse{
		Progress: StudyProgress{Finished: 0, Total: 0, StudyTime: 0},
	}
	points := MapStudyProgress(resp, time.Now())
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
	for _, p := range points {
		if p.MetricName == MetricTodayProgress && p.Value != 0 {
			t.Fatalf("expected 0 completion for zero total, got %v", p.Value)
		}
	}
}

func TestMapStudyProgress_Nil(t *testing.T) {
	if points := MapStudyProgress(nil, time.Now()); points != nil {
		t.Fatalf("expected nil, got %+v", points)
	}
}
