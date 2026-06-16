package mapper

import (
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestTimelogToProtoFormatsOptionalEndTime(t *testing.T) {
	endTime := time.Date(2026, 5, 30, 12, 30, 0, 0, time.UTC)
	log := &domain.TimeLog{
		ID:         10,
		UserID:     1,
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		EndTime:    &endTime,
		CategoryID: 3,
		TaskID:     &[]int32{20}[0],
		Remark:     "implementation",
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := TimelogToProto(log)
	if got.GetId() != 10 || got.GetStartTime() != "2026-05-30T12:00:00Z" || got.GetEndTime() != "2026-05-30T12:30:00Z" {
		t.Fatalf("unexpected mapped timelog: %#v", got)
	}
	if got.GetCategoryId() != 3 || got.GetTaskId() != 20 || got.GetRemark() != "implementation" {
		t.Fatalf("unexpected mapped timelog fields: %#v", got)
	}
}

func TestTimelogToProtoNil(t *testing.T) {
	if TimelogToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestTimelogsToProto(t *testing.T) {
	logs := []domain.TimeLog{
		{ID: 1, StartTime: time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC), CategoryID: 1},
		{ID: 2, StartTime: time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC), CategoryID: 1},
	}
	got := TimelogsToProto(logs)
	if len(got) != 2 || got[0].GetId() != 1 || got[1].GetId() != 2 {
		t.Fatalf("unexpected mapped timelogs: %#v", got)
	}
}

func TestTimelogFromCreateRequest(t *testing.T) {
	start := "2026-05-30T12:00:00Z"
	end := "2026-05-30T13:00:00Z"
	req := &timelogv1.CreateTimelogRequest{
		StartTime:  start,
		EndTime:    &end,
		CategoryId: 3,
		TaskId:     &[]int32{5}[0],
		Remark:     "note",
	}
	got, err := TimelogFromCreateRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.StartTime.IsZero() || got.CategoryID != 3 || got.Remark != "note" {
		t.Fatalf("unexpected domain timelog: %#v", got)
	}
	if got.EndTime == nil || got.TaskID == nil || *got.TaskID != 5 {
		t.Fatalf("expected end time and task id: %#v", got)
	}
}

func TestTimelogFromCreateRequestNil(t *testing.T) {
	got, err := TimelogFromCreateRequest(nil)
	if err != nil || got != nil {
		t.Fatalf("expected nil, got (%v, %v)", got, err)
	}
}

func TestApplyTimelogUpdate(t *testing.T) {
	tl := &domain.TimeLog{
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		CategoryID: 1,
		Remark:     "old",
	}
	newStart := "2026-05-30T14:00:00Z"
	newEnd := "2026-05-30T15:00:00Z"
	newRemark := "new"
	newCategoryID := int32(2)
	newTaskID := int32(7)
	req := &timelogv1.UpdateTimelogRequest{
		StartTime:  &newStart,
		EndTime:    &newEnd,
		CategoryId: &newCategoryID,
		TaskId:     &newTaskID,
		Remark:     &newRemark,
	}
	if err := ApplyTimelogUpdate(tl, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl.CategoryID != 2 || tl.Remark != "new" || tl.TaskID == nil || *tl.TaskID != 7 {
		t.Fatalf("unexpected updated timelog: %#v", tl)
	}
}

func TestApplyTimelogUpdateNil(t *testing.T) {
	if err := ApplyTimelogUpdate(nil, nil); err != nil {
		t.Fatalf("expected no error for nil inputs, got %v", err)
	}
}
