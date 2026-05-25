package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var singaporeLocation *time.Location

func init() {
	var err error
	singaporeLocation, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		// Fallback to UTC+8 if timezone data is not available
		singaporeLocation = time.FixedZone("SGT", 8*60*60)
	}
}

// formatMCPResponse wraps the response in the standard format to prevent LLM hallucinations
func formatMCPResponse(summaryText string, data interface{}) (*mcp.CallToolResult, interface{}, error) {
	// Add summary to data
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("data must be a map[string]interface{}")
	}
	dataMap["_summary"] = summaryText

	jsonBytes, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	responseText := string(jsonBytes)

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{
			Text: responseText,
		}},
	}, nil, nil
}

// Tool handlers with correct MCP signature
type DateInfoParams struct{}

func GetDateInfo(ctx context.Context, req *mcp.CallToolRequest, args DateInfoParams) (*mcp.CallToolResult, interface{}, error) {
	now := time.Now().In(singaporeLocation)
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	weekday := now.Weekday()
	// 以周一为一周的开始
	daysSinceMonday := (int(weekday) + 6) % 7
	monday := now.AddDate(0, 0, -daysSinceMonday)
	sunday := monday.AddDate(0, 0, 6)
	weekRange := []string{
		monday.Format("2006-01-02"),
		sunday.Format("2006-01-02"),
	}

	response := map[string]interface{}{
		"timezone":   "Asia/Singapore (SGT, UTC+8)",
		"now":        now.Format("2006-01-02 15:04:05"),
		"today":      today,
		"yesterday":  yesterday,
		"weekday":    weekday.String(),
		"week_range": weekRange,
	}
	summaryText := "当前日期和时间信息，包括今天、昨天和本周日期范围"
	return formatMCPResponse(summaryText, response)
}

func GetTimeLogsByDateRange(ctx context.Context, req *mcp.CallToolRequest, args DateRangeParams) (*mcp.CallToolResult, interface{}, error) {
	startDateStr := args.StartDate
	endDateStr := args.EndDate

	// 使用 model 层的函数，自动处理时区转换
	timeLogs, err := model.ListTimeLogsByLocalDateRange(server.db, startDateStr, endDateStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get time logs by date range: %w", err)
	}

	// 获取新加坡时区用于格式化输出
	sgLocation := model.GetSingaporeLocation()

	var result []map[string]interface{}
	totalDuration := time.Duration(0)

	for _, tl := range timeLogs {
		duration := time.Duration(0)
		durationStr := ""

		if tl.EndTime != nil {
			duration = tl.EndTime.Sub(tl.StartTime)
			totalDuration += duration
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			durationStr = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			durationStr = "ongoing"
		}

		entry := map[string]interface{}{
			"id":         tl.ID,
			"start_time": tl.StartTime.In(sgLocation).Format("2006-01-02 15:04:05"),
			"end_time":   nil,
			"duration":   durationStr,
			"remarks":    tl.Remark,
		}

		if tl.EndTime != nil {
			entry["end_time"] = tl.EndTime.In(sgLocation).Format("2006-01-02 15:04:05")
		}

		result = append(result, entry)
	}

	totalHours := int(totalDuration.Hours())
	totalMinutes := int(totalDuration.Minutes()) % 60

	response := map[string]interface{}{
		"time_logs":      result,
		"count":          len(result),
		"date_range":     fmt.Sprintf("%s to %s", startDateStr, endDateStr[:10]),
		"total_duration": fmt.Sprintf("%dh %dm", totalHours, totalMinutes),
	}

	summaryText := fmt.Sprintf("Found %d time logs from %s to %s, total duration: %dh %dm", len(result), startDateStr, endDateStr[:10], totalHours, totalMinutes)
	return formatMCPResponse(summaryText, response)
}

func GetTasksByStatus(ctx context.Context, req *mcp.CallToolRequest, args TaskStatusParams) (*mcp.CallToolResult, interface{}, error) {
	statusStr := args.Status

	// Include all tasks (suspended and completed) to filter by status in application code
	tasks, err := model.GetAllTasks(server.db, true, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		isCompleted := task.IsCompleted != nil && *task.IsCompleted

		// Filter by status
		if statusStr == "completed" && !isCompleted {
			continue
		}
		if statusStr == "pending" && isCompleted {
			continue
		}

		categoryName := ""
		categoryColor := ""
		if task.CategoryID > 0 {
			if cat, err := model.GetCategoryByID(server.db, int32(task.CategoryID)); err == nil && cat != nil {
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
			"due_date":          task.DueDate.In(singaporeLocation).Format("2006-01-02"),
			"estimated_minutes": task.EstimatedMinutes,
			"is_completed":      isCompleted,
			"created_at":        task.CreatedAt.In(singaporeLocation).Format("2006-01-02 15:04:05"),
		}

		if task.CompletedAt != nil {
			entry["completed_at"] = task.CompletedAt.In(singaporeLocation).Format("2006-01-02 15:04:05")
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"tasks":  result,
		"count":  len(result),
		"status": statusStr,
	}

	summaryText := fmt.Sprintf("Found %d %s tasks", len(result), statusStr)
	return formatMCPResponse(summaryText, response)
}

func GetCurrentActivity(ctx context.Context, req *mcp.CallToolRequest, args CurrentActivityParams) (*mcp.CallToolResult, interface{}, error) {
	timeLogs, err := model.ListTimeLogsWithOptions(server.db, 0, "start_time DESC", "end_time IS NULL")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current activity: %w", err)
	}

	var result []map[string]interface{}
	for _, tl := range timeLogs {
		duration := time.Since(tl.StartTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		entry := map[string]interface{}{
			"id":         tl.ID,
			"start_time": tl.StartTime.In(singaporeLocation).Format("2006-01-02 15:04:05"),
			"duration":   fmt.Sprintf("%dh %dm", hours, minutes),
			"remarks":    tl.Remark,
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"active_logs": result,
		"count":       len(result),
	}

	summaryText := fmt.Sprintf("Found %d active time logs", len(result))
	return formatMCPResponse(summaryText, response)
}

func GetActiveConstraints(ctx context.Context, req *mcp.CallToolRequest, args ConstraintParams) (*mcp.CallToolResult, interface{}, error) {
	constraints, err := model.GetActiveConstraints(server.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active constraints: %w", err)
	}

	var result []map[string]interface{}
	for _, constraint := range constraints {
		isActive := constraint.IsActive != nil && *constraint.IsActive
		createdAt := ""
		if constraint.CreatedAt != nil {
			createdAt = constraint.CreatedAt.In(singaporeLocation).Format("2006-01-02 15:04:05")
		}

		entry := map[string]interface{}{
			"id":               constraint.ID,
			"description":      constraint.Description,
			"punishment_quote": constraint.PunishmentQuote,
			"start_date":       constraint.StartDate.In(singaporeLocation).Format("2006-01-02"),
			"is_active":        isActive,
			"created_at":       createdAt,
		}

		if constraint.EndDate != nil {
			entry["end_date"] = constraint.EndDate.In(singaporeLocation).Format("2006-01-02")
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

	summaryText := fmt.Sprintf("Found %d active constraints", len(result))
	return formatMCPResponse(summaryText, response)
}
