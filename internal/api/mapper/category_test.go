package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

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
