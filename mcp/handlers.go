package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func formatMCPResponse(summaryText string, data map[string]interface{}) (*mcp.CallToolResult, interface{}, error) {
	data["_summary"] = summaryText

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{
			Text: string(jsonBytes),
		}},
	}, nil, nil
}

type DateInfoParams struct{}

func GetDateInfo(ctx context.Context, req *mcp.CallToolRequest, args DateInfoParams) (*mcp.CallToolResult, interface{}, error) {
	loc := model.GetSingaporeLocation()
	now := time.Now().In(loc)
	weekday := now.Weekday()
	daysSinceMonday := (int(weekday) + 6) % 7
	monday := now.AddDate(0, 0, -daysSinceMonday)
	sunday := monday.AddDate(0, 0, 6)

	response := map[string]interface{}{
		"timezone":   "Asia/Singapore (SGT, UTC+8)",
		"now":        service.FormatSGDateTime(now),
		"today":      service.TodaySGDateString(),
		"yesterday":  now.AddDate(0, 0, -1).Format("2006-01-02"),
		"weekday":    weekday.String(),
		"week_range": []string{monday.Format("2006-01-02"), sunday.Format("2006-01-02")},
	}
	return formatMCPResponse("当前日期和时间信息，包括今天、昨天和本周日期范围", response)
}

func GetTimeLogsByDateRange(ctx context.Context, req *mcp.CallToolRequest, args DateRangeParams) (*mcp.CallToolResult, interface{}, error) {
	if args.ActiveOnly {
		timeLogs, err := service.ListActiveTimeLogs()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get current activity: %w", err)
		}

		var result []map[string]interface{}
		for _, tl := range timeLogs {
			entry := map[string]interface{}{
				"id":         tl.ID,
				"start_time": service.FormatSGDateTime(tl.StartTime),
				"duration":   service.FormatDuration(time.Since(tl.StartTime)),
				"remark":     tl.Remark,
			}
			result = append(result, entry)
		}

		response := map[string]interface{}{
			"active_logs": result,
			"count":       len(result),
		}
		return formatMCPResponse(fmt.Sprintf("Found %d active time logs", len(result)), response)
	}

	startDate := args.StartDate
	endDate := args.EndDate
	if startDate == "" {
		startDate = service.TodaySGDateString()
	}
	if endDate == "" {
		endDate = startDate
	}

	timeLogs, err := service.ListTimeLogsByLocalDateRange(startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get time logs by date range: %w", err)
	}

	var result []map[string]interface{}
	totalDuration := time.Duration(0)

	for _, tl := range timeLogs {
		duration := time.Duration(0)
		durationStr := "ongoing"

		if tl.EndTime != nil {
			duration = tl.EndTime.Sub(tl.StartTime)
			totalDuration += duration
			durationStr = service.FormatDuration(duration)
		}

		entry := map[string]interface{}{
			"id":         tl.ID,
			"start_time": service.FormatSGDateTime(tl.StartTime),
			"end_time":   nil,
			"duration":   durationStr,
			"remark":     tl.Remark,
		}

		if tl.EndTime != nil {
			entry["end_time"] = service.FormatSGDateTime(*tl.EndTime)
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"time_logs":      result,
		"count":          len(result),
		"date_range":     fmt.Sprintf("%s to %s", startDate, endDate),
		"total_duration": service.FormatDuration(totalDuration),
	}

	summaryText := fmt.Sprintf("Found %d time logs from %s to %s, total duration: %s", len(result), startDate, endDate, service.FormatDuration(totalDuration))
	return formatMCPResponse(summaryText, response)
}

func GetTasksByStatus(ctx context.Context, req *mcp.CallToolRequest, args TaskStatusParams) (*mcp.CallToolResult, interface{}, error) {
	tasks, err := service.ListTasksByCompletionStatus(args.Status)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		isCompleted := task.IsCompleted != nil && *task.IsCompleted

		categoryName := ""
		categoryColor := ""
		if task.CategoryID > 0 {
			if cat, err := service.GetCategoryByID(int32(task.CategoryID)); err == nil && cat != nil {
				categoryName = cat.Name
				if cat.Color != nil {
					categoryColor = *cat.Color
				}
			}
		}

		entry := map[string]interface{}{
			"id":                task.ID,
			"title":             task.Title,
			"description":       task.Description,
			"category":          categoryName,
			"category_color":    categoryColor,
			"due_date":          service.FormatSGDate(task.DueDate),
			"estimated_minutes": task.EstimatedMinutes,
			"is_completed":      isCompleted,
			"created_at":        service.FormatSGDateTimePtr(task.CreatedAt),
		}

		if task.CompletedAt != nil {
			entry["completed_at"] = service.FormatSGDateTime(*task.CompletedAt)
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"tasks":  result,
		"count":  len(result),
		"status": args.Status,
	}

	return formatMCPResponse(fmt.Sprintf("Found %d %s tasks", len(result), args.Status), response)
}

func GetActiveConstraints(ctx context.Context, req *mcp.CallToolRequest, args ConstraintParams) (*mcp.CallToolResult, interface{}, error) {
	constraints, err := service.GetActiveConstraints()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active constraints: %w", err)
	}

	var result []map[string]interface{}
	for _, constraint := range constraints {
		isActive := constraint.IsActive != nil && *constraint.IsActive

		entry := map[string]interface{}{
			"id":               constraint.ID,
			"description":      constraint.Description,
			"punishment_quote": constraint.PunishmentQuote,
			"start_date":       service.FormatSGDate(constraint.StartDate),
			"is_active":        isActive,
			"created_at":       service.FormatSGDateTimePtr(constraint.CreatedAt),
		}

		if constraint.EndDate != nil {
			entry["end_date"] = service.FormatSGDate(*constraint.EndDate)
		}
		if constraint.EndReason != nil && *constraint.EndReason != "" {
			entry["end_reason"] = *constraint.EndReason
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"constraints": result,
		"count":       len(result),
	}

	return formatMCPResponse(fmt.Sprintf("Found %d active constraints", len(result)), response)
}

func ListCategories(ctx context.Context, req *mcp.CallToolRequest, args CategoryListParams) (*mcp.CallToolResult, interface{}, error) {
	categories, err := service.ListCategories()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list categories: %w", err)
	}

	var result []map[string]interface{}
	for _, cat := range categories {
		entry := map[string]interface{}{
			"id":    cat.ID,
			"name":  cat.Name,
			"level": int32(0),
		}
		if cat.Level != nil {
			entry["level"] = *cat.Level
		}
		if cat.Color != nil {
			entry["color"] = *cat.Color
		}
		if cat.ParentID != nil && *cat.ParentID > 0 {
			entry["parent_id"] = *cat.ParentID
		}
		result = append(result, entry)
	}

	response := map[string]interface{}{
		"categories": result,
		"count":      len(result),
	}

	return formatMCPResponse(fmt.Sprintf("Found %d categories", len(result)), response)
}

func CreateTimeLog(ctx context.Context, req *mcp.CallToolRequest, args CreateTimeLogParams) (*mcp.CallToolResult, interface{}, error) {
	createdLog, err := service.CreateTimeLogFromMCPInput(service.CreateTimeLogMCPInput{
		CategoryID: args.CategoryID,
		StartTime:  args.StartTime,
		EndTime:    args.EndTime,
		TaskID:     args.TaskID,
		Remark:     args.Remark,
	})
	if err != nil {
		return nil, nil, err
	}

	response := map[string]interface{}{
		"id":          *createdLog.ID,
		"category_id": createdLog.CategoryID,
		"start_time":  service.FormatSGDateTime(createdLog.StartTime),
		"remark":      createdLog.Remark,
	}
	if createdLog.EndTime != nil {
		response["end_time"] = service.FormatSGDateTime(*createdLog.EndTime)
	}
	if createdLog.TaskID != nil {
		response["task_id"] = *createdLog.TaskID
	}

	return formatMCPResponse("Time log created successfully", response)
}

func UpdateTimeLog(ctx context.Context, req *mcp.CallToolRequest, args UpdateTimeLogParams) (*mcp.CallToolResult, interface{}, error) {
	updatedLog, err := service.UpdateTimeLogFromMCPInput(service.UpdateTimeLogMCPInput{
		ID:         args.ID,
		CategoryID: args.CategoryID,
		StartTime:  args.StartTime,
		EndTime:    args.EndTime,
		TaskID:     args.TaskID,
		Remark:     args.Remark,
	})
	if err != nil {
		return nil, nil, err
	}

	response := map[string]interface{}{
		"id":          *updatedLog.ID,
		"category_id": updatedLog.CategoryID,
		"start_time":  service.FormatSGDateTime(updatedLog.StartTime),
		"remark":      updatedLog.Remark,
	}
	if updatedLog.EndTime != nil {
		response["end_time"] = service.FormatSGDateTime(*updatedLog.EndTime)
	}
	if updatedLog.TaskID != nil {
		response["task_id"] = *updatedLog.TaskID
	}

	return formatMCPResponse("Time log updated successfully", response)
}
