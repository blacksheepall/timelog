package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCreateTaskFromMCPInput(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	task, err := svc.CreateTaskFromMCPInput(CreateTaskMCPInput{
		Title:            "Review plan",
		Description:      "Read and approve the implementation plan",
		CategoryID:       cat.ID,
		EstimatedMinutes: 30,
	})
	if err != nil {
		t.Fatalf("CreateTaskFromMCPInput: %v", err)
	}
	if task.ID == 0 {
		t.Fatal("expected task ID to be set")
	}
	if task.Title != "Review plan" {
		t.Fatalf("expected title %q, got %q", "Review plan", task.Title)
	}
	if task.CategoryID != cat.ID {
		t.Fatalf("expected category_id %d, got %d", cat.ID, task.CategoryID)
	}
	if task.EstimatedMinutes != 30 {
		t.Fatalf("expected estimated_minutes 30, got %d", task.EstimatedMinutes)
	}
}

func TestCreateTaskFromMCPInputMissingTitle(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.CreateTaskFromMCPInput(CreateTaskMCPInput{})
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestCreateTaskFromMCPInputMissingCategory(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.CreateTaskFromMCPInput(CreateTaskMCPInput{
		Title: "No category",
	})
	if err == nil {
		t.Fatal("expected error for missing category")
	}
}

func TestCreateTaskFromMCPInputInvalidDueDate(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	_, err := svc.CreateTaskFromMCPInput(CreateTaskMCPInput{
		Title:      "Bad date",
		CategoryID: cat.ID,
		DueDate:    "bad",
	})
	if err == nil {
		t.Fatal("expected error for invalid due_date")
	}
}

func TestCreateTaskFromMCPInputInvalidCategory(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.CreateTaskFromMCPInput(CreateTaskMCPInput{
		Title:      "Orphan",
		CategoryID: 9999,
	})
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestUpdateTaskFromMCPInput(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	newTitle := "updated"
	task, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:    created.ID,
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("UpdateTaskFromMCPInput: %v", err)
	}
	if task.Title != "updated" {
		t.Fatalf("expected title updated, got %q", task.Title)
	}
}

func TestUpdateTaskFromMCPInputNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	newTitle := "updated"
	_, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:    9999,
		Title: &newTitle,
	})
	if err == nil {
		t.Fatal("expected error for missing task")
	}
}

func TestUpdateTaskFromMCPInputInvalidCategory(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	_, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:         created.ID,
		CategoryID: int32Ptr(9999),
	})
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestUpdateTaskFromMCPInputZeroCategory(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	_, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:         created.ID,
		CategoryID: int32Ptr(0),
	})
	if err == nil {
		t.Fatal("expected error for zero category")
	}
}

func TestUpdateTaskFromMCPInputEmptyTitle(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	_, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:    created.ID,
		Title: strPtr(""),
	})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestUpdateTaskFromMCPInputInvalidDueDate(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	_, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:      created.ID,
		DueDate: strPtr("bad"),
	})
	if err == nil {
		t.Fatal("expected error for invalid due_date")
	}
}

func TestUpdateTaskFromMCPInputPartialFields(t *testing.T) {
	svc, _ := newTestService(t)
	cat := &domain.Category{Name: "mcp-task-test"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("create category: %v", err)
	}

	created := &domain.Task{Title: "orig", CategoryID: cat.ID}
	if err := svc.CreateTask(created); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	description := "updated description"
	dueDate := "2025-12-31"
	estimated := int32(60)
	task, err := svc.UpdateTaskFromMCPInput(UpdateTaskMCPInput{
		ID:               created.ID,
		Description:      &description,
		DueDate:          &dueDate,
		EstimatedMinutes: &estimated,
	})
	if err != nil {
		t.Fatalf("UpdateTaskFromMCPInput: %v", err)
	}
	if task.Description != description {
		t.Fatalf("expected description %q, got %q", description, task.Description)
	}
	expectedDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	if !task.DueDate.Equal(expectedDate) {
		t.Fatalf("expected due_date %v, got %v", expectedDate, task.DueDate)
	}
	if task.EstimatedMinutes != estimated {
		t.Fatalf("expected estimated_minutes %d, got %d", estimated, task.EstimatedMinutes)
	}
}
