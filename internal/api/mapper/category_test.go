package mapper

import (
	"errors"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/core/errs"
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
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

func TestCategoryFromCreateRequestRejectsClientSortOrder(t *testing.T) {
	sortOrder := int32(1)
	req := &timelogv1.CreateCategoryRequest{
		Name:      "Root",
		SortOrder: &sortOrder,
	}
	_, err := CategoryFromCreateRequest(req)
	if err == nil {
		t.Fatal("expected error for client-provided sort_order")
	}
}

func TestCategoryFromCreateRequestHappyPath(t *testing.T) {
	parentID := int32(3)
	req := &timelogv1.CreateCategoryRequest{
		Name:        "Work",
		Color:       "#3366ff",
		Description: "Work category",
		ParentId:    &parentID,
	}
	got, err := CategoryFromCreateRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Work" || got.Color != "#3366ff" || got.Description != "Work category" || got.ParentID == nil || *got.ParentID != 3 {
		t.Fatalf("unexpected domain category: %#v", got)
	}
}

func TestCategoryToProtoKeepsSnakeCaseFields(t *testing.T) {
	parentID := int32(3)
	c := &domain.Category{
		ID:          7,
		Name:        "Coding",
		Color:       "#3366ff",
		Description: "Deep work",
		ParentID:    &parentID,
		Level:       1,
		SortOrder:   2,
		Path:        "/3",
		CreatedAt:   time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := CategoryToProto(c)
	if got.GetName() != "Coding" || got.GetParentId() != 3 || got.GetCreatedAt() != "2026-05-30T12:00:00Z" {
		t.Fatalf("unexpected mapped category: %#v", got)
	}
}

func TestCategoryToProtoNil(t *testing.T) {
	if CategoryToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestCategoriesToProto(t *testing.T) {
	cats := []domain.Category{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	}
	got := CategoriesToProto(cats)
	if len(got) != 2 || got[0].GetName() != "A" || got[1].GetName() != "B" {
		t.Fatalf("unexpected mapped categories: %#v", got)
	}
}

func TestCategoryTreeNodeToProto(t *testing.T) {
	node := &domain.CategoryNode{
		Category: domain.Category{ID: 1, Name: "Root"},
		Children: []*domain.CategoryNode{
			{Category: domain.Category{ID: 2, Name: "Child"}},
		},
	}
	got := CategoryTreeNodeToProto(node)
	if got.GetCategory().GetName() != "Root" || len(got.GetChildren()) != 1 || got.GetChildren()[0].GetCategory().GetName() != "Child" {
		t.Fatalf("unexpected mapped tree node: %#v", got)
	}
}

func TestCategoryTreeToProto(t *testing.T) {
	nodes := []*domain.CategoryNode{
		{Category: domain.Category{ID: 1, Name: "Root"}},
	}
	got := CategoryTreeToProto(nodes)
	if len(got) != 1 || got[0].GetCategory().GetName() != "Root" {
		t.Fatalf("unexpected mapped tree: %#v", got)
	}
}

func TestValidateMoveCategoryRequest(t *testing.T) {
	if err := ValidateMoveCategoryRequest(nil); err == nil {
		t.Fatal("expected error for nil request")
	}
	parentZero := int32(0)
	req := &timelogv1.MoveCategoryRequest{ParentId: &parentZero}
	if !errors.Is(ValidateMoveCategoryRequest(req), errs.ErrInvalidParentID) {
		t.Fatalf("expected ErrInvalidParentID")
	}
	parentID := int32(3)
	req2 := &timelogv1.MoveCategoryRequest{ParentId: &parentID}
	if err := ValidateMoveCategoryRequest(req2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyCategoryUpdate(t *testing.T) {
	cat := &domain.Category{Name: "Old", Color: "#000000", Description: "OldDesc"}
	newName := "New"
	newColor := "#FFFFFF"
	newDesc := "NewDesc"
	req := &timelogv1.UpdateCategoryRequest{
		Name:        &newName,
		Color:       &newColor,
		Description: &newDesc,
	}
	if err := ApplyCategoryUpdate(cat, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.Name != "New" || cat.Color != "#FFFFFF" || cat.Description != "NewDesc" {
		t.Fatalf("unexpected updated category: %#v", cat)
	}
}

func TestApplyCategoryUpdateRejectsParentID(t *testing.T) {
	cat := &domain.Category{Name: "Old"}
	parentID := int32(1)
	req := &timelogv1.UpdateCategoryRequest{ParentId: &parentID}
	if !errors.Is(ApplyCategoryUpdate(cat, req), errs.ErrInvalidParentID) {
		t.Fatalf("expected ErrInvalidParentID")
	}
}
