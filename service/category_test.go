package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCreateCategoryRejectsInvalidParentID(t *testing.T) {
	svc, _ := newTestService(t)
	parentID := int32(0)
	if err := svc.CreateCategory(&domain.Category{Name: "bad", ParentID: &parentID}); err == nil {
		t.Fatal("expected error for invalid parent_id")
	}
}

func TestCreateCategoryRejectsExceedingMaxLevel(t *testing.T) {
	svc, _ := newTestService(t)
	root := &domain.Category{Name: "Root"}
	if err := svc.CreateCategory(root); err != nil {
		t.Fatalf("CreateCategory root: %v", err)
	}
	child := &domain.Category{Name: "Child", ParentID: &root.ID}
	if err := svc.CreateCategory(child); err != nil {
		t.Fatalf("CreateCategory child: %v", err)
	}
	grandchild := &domain.Category{Name: "Grandchild", ParentID: &child.ID}
	if err := svc.CreateCategory(grandchild); err != nil {
		t.Fatalf("CreateCategory grandchild: %v", err)
	}
	fourth := &domain.Category{Name: "TooDeep", ParentID: &grandchild.ID}
	if err := svc.CreateCategory(fourth); err == nil {
		t.Fatal("expected error exceeding max level")
	}
}

func TestGetCategoryByIDNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.GetCategoryByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent category")
	}
}

func TestGetCategoryByNameNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.GetCategoryByName("missing", nil)
	if err == nil {
		t.Fatal("expected error for non-existent category name")
	}
}

func TestUpdateCategoryNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.UpdateCategory(&domain.Category{ID: 9999, Name: "ghost"}); err == nil {
		t.Fatal("expected error updating non-existent category")
	}
}

func TestMoveCategoryRejectsSelf(t *testing.T) {
	svc, _ := newTestService(t)
	root := &domain.Category{Name: "Root"}
	if err := svc.CreateCategory(root); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	parentID := root.ID
	if err := svc.MoveCategory(root.ID, &parentID); err == nil {
		t.Fatal("expected error moving category to itself")
	}
}

func TestMoveCategoryRejectsDescendant(t *testing.T) {
	svc, _ := newTestService(t)
	root := &domain.Category{Name: "Root"}
	if err := svc.CreateCategory(root); err != nil {
		t.Fatalf("CreateCategory root: %v", err)
	}
	child := &domain.Category{Name: "Child", ParentID: &root.ID}
	if err := svc.CreateCategory(child); err != nil {
		t.Fatalf("CreateCategory child: %v", err)
	}
	if err := svc.MoveCategory(root.ID, &child.ID); err == nil {
		t.Fatal("expected error moving category to its own child")
	}
}
