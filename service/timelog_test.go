package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/core/errs"
)

func TestCreateTimeLogRejectsOngoingConflict(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("CreateTimeLog first: %v", err)
	}
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); !errors.Is(err, errs.ErrOngoingTimeLogExists) {
		t.Fatalf("expected ErrOngoingTimeLogExists, got %v", err)
	}
}

func TestGetTimeLogByIDNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.GetTimeLogByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent timelog")
	}
}

func TestUpdateTimeLogNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.UpdateTimeLog(context.Background(), &domain.TimeLog{ID: 9999, StartTime: time.Now().UTC(), CategoryID: 1}); err == nil {
		t.Fatal("expected error updating non-existent timelog")
	}
}

func TestDeleteTimeLogNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.DeleteTimeLog(context.Background(), 9999); err == nil {
		t.Fatal("expected error deleting non-existent timelog")
	}
}

func TestCreateTimeLogFromMCPInputRejectsOngoingConflict(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	if _, err := svc.CreateTimeLogFromMCPInput(context.Background(), CreateTimeLogMCPInput{CategoryID: 1}); err != nil {
		t.Fatalf("CreateTimeLogFromMCPInput first: %v", err)
	}
	if _, err := svc.CreateTimeLogFromMCPInput(context.Background(), CreateTimeLogMCPInput{CategoryID: 1}); !errors.Is(err, errs.ErrOngoingTimeLogExists) {
		t.Fatalf("expected ErrOngoingTimeLogExists, got %v", err)
	}
}

func TestCreateTimeLogFromMCPInputRejectsInvalidEndTime(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	if _, err := svc.CreateTimeLogFromMCPInput(context.Background(), CreateTimeLogMCPInput{CategoryID: 1, EndTime: "bad"}); err == nil {
		t.Fatal("expected error for invalid end time")
	}
}

func TestUpdateTimeLogFromMCPInputNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.UpdateTimeLogFromMCPInput(context.Background(), UpdateTimeLogMCPInput{ID: 9999}); err == nil {
		t.Fatal("expected error for non-existent timelog")
	}
}

func TestUpdateTimeLogFromMCPInputRejectsInvalidCategory(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	created, err := svc.CreateTimeLogFromMCPInput(context.Background(), CreateTimeLogMCPInput{CategoryID: 1})
	if err != nil {
		t.Fatalf("CreateTimeLogFromMCPInput: %v", err)
	}

	if _, err := svc.UpdateTimeLogFromMCPInput(context.Background(), UpdateTimeLogMCPInput{ID: created.ID, CategoryID: 999}); err == nil {
		t.Fatal("expected error for invalid category")
	}
}
