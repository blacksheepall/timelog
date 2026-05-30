package mapper

import (
	"errors"
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/errs"
)

func TestCategoryFromCreateRequestRejectsInvalidParentID(t *testing.T) {
	parentZero := int32(0)
	req := &timelogv1.CreateCategoryRequest{
		Name:        "Root",
		Color:       "#3366ff",
		Description: "Root category",
		ParentId:    &parentZero,
	}
	_, err := CategoryFromCreateRequest(req)
	if !errors.Is(err, errs.ErrInvalidParentID) {
		t.Fatalf("expected ErrInvalidParentID, got %v", err)
	}
}

func TestCategoryFromCreateRequestRejectsClientLevel(t *testing.T) {
	level := int32(1)
	req := &timelogv1.CreateCategoryRequest{
		Name:        "Root",
		Color:       "#3366ff",
		Description: "Root category",
		Level:       &level,
	}
	_, err := CategoryFromCreateRequest(req)
	if err == nil {
		t.Fatal("expected error for client-provided level")
	}
}

func TestCategoryToProtoKeepsSnakeCaseFields(t *testing.T) {
	id := int32(7)
	parentID := int32(3)
	level := int32(1)
	sortOrder := int32(2)
	color := "#3366ff"
	description := "Deep work"
	path := "/3"
	c := &gen.Category{
		ID:          &id,
		Name:        "Coding",
		Color:       &color,
		Description: &description,
		ParentID:    &parentID,
		Level:       &level,
		SortOrder:   &sortOrder,
		Path:        &path,
		CreatedAt:   time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := CategoryToProto(c)
	if got.GetName() != "Coding" || got.GetParentId() != 3 || got.GetCreatedAt() != "2026-05-30T12:00:00Z" {
		t.Fatalf("unexpected mapped category: %#v", got)
	}
}
