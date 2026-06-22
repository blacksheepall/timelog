package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

type CreateTaskMCPInput struct {
	Title            string
	Description      string
	CategoryID       int32
	DueDate          string
	EstimatedMinutes int32
}

type UpdateTaskMCPInput struct {
	ID               int32
	Title            *string
	Description      *string
	CategoryID       *int32
	DueDate          *string
	EstimatedMinutes *int32
}

func (s *Service) CreateTaskFromMCPInput(input CreateTaskMCPInput) (*domain.Task, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if input.CategoryID <= 0 {
		return nil, fmt.Errorf("category_id is required")
	}
	if _, err := s.GetCategoryByID(input.CategoryID); err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	dueDate, err := parseTaskDueDate(input.DueDate)
	if err != nil {
		return nil, err
	}

	task := &domain.Task{
		Title:            input.Title,
		Description:      input.Description,
		CategoryID:       input.CategoryID,
		DueDate:          dueDate,
		EstimatedMinutes: input.EstimatedMinutes,
	}

	if err := s.CreateTask(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return s.GetTaskByID(task.ID)
}

func (s *Service) UpdateTaskFromMCPInput(input UpdateTaskMCPInput) (*domain.Task, error) {
	task, err := s.GetTaskByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	if input.Title != nil {
		if *input.Title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.CategoryID != nil {
		if *input.CategoryID <= 0 {
			return nil, fmt.Errorf("category_id must be a positive integer")
		}
		if _, err := s.GetCategoryByID(*input.CategoryID); err != nil {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		task.CategoryID = *input.CategoryID
	}
	if input.DueDate != nil {
		dueDate, err := parseTaskDueDate(*input.DueDate)
		if err != nil {
			return nil, err
		}
		task.DueDate = dueDate
	}
	if input.EstimatedMinutes != nil {
		task.EstimatedMinutes = *input.EstimatedMinutes
	}

	if err := s.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return s.GetTaskByID(input.ID)
}

func parseTaskDueDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid due_date format, expected YYYY-MM-DD: %w", err)
	}
	return t, nil
}
