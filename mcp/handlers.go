package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var singaporeLocation *time.Location

func init() {
	var err error
	singaporeLocation, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		singaporeLocation = time.FixedZone("SGT", 8*60*60)
	}
}

func formatSGDateTime(t time.Time) string {
	return t.In(singaporeLocation).Format("2006-01-02 15:04:05")
}

func formatSGDate(t time.Time) string {
	return t.In(singaporeLocation).Format("2006-01-02")
}

func formatSGDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(singaporeLocation).Format("2006-01-02 15:04:05")
}

func formatSGDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(singaporeLocation).Format("2006-01-02")
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

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

// Tool handlers with correct MCP signature
type DateInfoParams struct{}

func GetDateInfo(ctx context.Context, req *mcp.CallToolRequest, args DateInfoParams) (*mcp.CallToolResult, interface{}, error) {
	now := time.Now().In(singaporeLocation)
	weekday := now.Weekday()
	daysSinceMonday := (int(weekday) + 6) % 7
	monday := now.AddDate(0, 0, -daysSinceMonday)
	sunday := monday.AddDate(0, 0, 6)

	response := map[string]interface{}{
		"timezone":   "Asia/Singapore (SGT, UTC+8)",
		"now":        formatSGDateTime(now),
		"today":      now.Format("2006-01-02"),
		"yesterday":  now.AddDate(0, 0, -1).Format("2006-01-02"),
		"weekday":    weekday.String(),
		"week_range": []string{monday.Format("2006-01-02"), sunday.Format("2006-01-02")},
	}
	return formatMCPResponse("当前日期和时间信息，包括今天、昨天和本周日期范围", response)
}

func GetTimeLogsByDateRange(ctx context.Context, req *mcp.CallToolRequest, args DateRangeParams) (*mcp.CallToolResult, interface{}, error) {
	if args.ActiveOnly {
		timeLogs, err := model.ListTimeLogsWithOptions(server.db, 0, "start_time DESC", "end_time IS NULL")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get current activity: %w", err)
		}

		var result []map[string]interface{}
		for _, tl := range timeLogs {
			entry := map[string]interface{}{
				"id":         tl.ID,
				"start_time": formatSGDateTime(tl.StartTime),
				"duration":   formatDuration(time.Since(tl.StartTime)),
				"remarks":    tl.Remark,
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
		startDate = time.Now().In(singaporeLocation).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = startDate
	}

	timeLogs, err := model.ListTimeLogsByLocalDateRange(server.db, startDate, endDate)
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
			durationStr = formatDuration(duration)
		}

		entry := map[string]interface{}{
			"id":         tl.ID,
			"start_time": formatSGDateTime(tl.StartTime),
			"end_time":   nil,
			"duration":   durationStr,
			"remarks":    tl.Remark,
		}

		if tl.EndTime != nil {
			entry["end_time"] = formatSGDateTime(*tl.EndTime)
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"time_logs":      result,
		"count":          len(result),
		"date_range":     fmt.Sprintf("%s to %s", startDate, endDate),
		"total_duration": formatDuration(totalDuration),
	}

	summaryText := fmt.Sprintf("Found %d time logs from %s to %s, total duration: %s", len(result), startDate, endDate, formatDuration(totalDuration))
	return formatMCPResponse(summaryText, response)
}

func GetTasksByStatus(ctx context.Context, req *mcp.CallToolRequest, args TaskStatusParams) (*mcp.CallToolResult, interface{}, error) {
	tasks, err := model.GetAllTasks(server.db, true, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		isCompleted := task.IsCompleted != nil && *task.IsCompleted

		if args.Status == "completed" && !isCompleted {
			continue
		}
		if args.Status == "pending" && isCompleted {
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
			"due_date":          formatSGDate(task.DueDate),
			"estimated_minutes": task.EstimatedMinutes,
			"is_completed":      isCompleted,
			"created_at":        formatSGDateTimePtr(task.CreatedAt),
		}

		if task.CompletedAt != nil {
			entry["completed_at"] = formatSGDateTime(*task.CompletedAt)
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
	constraints, err := model.GetActiveConstraints(server.db)
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
			"start_date":       formatSGDate(constraint.StartDate),
			"is_active":        isActive,
			"created_at":       formatSGDateTimePtr(constraint.CreatedAt),
		}

		if constraint.EndDate != nil {
			entry["end_date"] = formatSGDate(*constraint.EndDate)
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
	categories, err := model.ListCategories(server.db)
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
	_, err := model.GetCategoryByID(server.db, args.CategoryID)
	if err != nil {
		return nil, nil, fmt.Errorf("category not found: %w", err)
	}

	tl := &gen.Timelog{
		CategoryID: args.CategoryID,
	}

	defaultUserID := int32(1)
	tl.UserID = &defaultUserID

	if args.StartTime != "" {
		st, err := time.ParseInLocation("2006-01-02 15:04:05", args.StartTime, singaporeLocation)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		tl.StartTime = st.UTC()
	} else {
		tl.StartTime = time.Now().UTC()
	}

	if args.EndTime != "" {
		et, err := time.ParseInLocation("2006-01-02 15:04:05", args.EndTime, singaporeLocation)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		etUTC := et.UTC()
		tl.EndTime = &etUTC
	}

	if args.TaskID > 0 {
		tl.TaskID = &args.TaskID
	}

	if args.Remark != "" {
		tl.Remark = &args.Remark
	}

	if err := model.CreateTimeLog(server.db, tl); err != nil {
		return nil, nil, fmt.Errorf("failed to create time log: %w", err)
	}

	createdLog, err := model.GetTimeLogByID(server.db, *tl.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve created time log: %w", err)
	}

	response := map[string]interface{}{
		"id":          *createdLog.ID,
		"category_id": createdLog.CategoryID,
		"start_time":  formatSGDateTime(createdLog.StartTime),
		"remark":      createdLog.Remark,
	}
	if createdLog.EndTime != nil {
		response["end_time"] = formatSGDateTime(*createdLog.EndTime)
	}
	if createdLog.TaskID != nil {
		response["task_id"] = *createdLog.TaskID
	}

	return formatMCPResponse("Time log created successfully", response)
}

func UpdateTimeLog(ctx context.Context, req *mcp.CallToolRequest, args UpdateTimeLogParams) (*mcp.CallToolResult, interface{}, error) {
	tl, err := model.GetTimeLogByID(server.db, args.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("time log not found: %w", err)
	}

	if args.CategoryID > 0 {
		_, err := model.GetCategoryByID(server.db, args.CategoryID)
		if err != nil {
			return nil, nil, fmt.Errorf("category not found: %w", err)
		}
		tl.CategoryID = args.CategoryID
	}

	if args.StartTime != "" {
		st, err := time.ParseInLocation("2006-01-02 15:04:05", args.StartTime, singaporeLocation)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		tl.StartTime = st.UTC()
	}

	if args.EndTime != "" {
		et, err := time.ParseInLocation("2006-01-02 15:04:05", args.EndTime, singaporeLocation)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		etUTC := et.UTC()
		tl.EndTime = &etUTC
	}

	if args.TaskID > 0 {
		tl.TaskID = &args.TaskID
	}

	if args.Remark != "" {
		tl.Remark = &args.Remark
	}

	if err := model.UpdateTimeLog(server.db, tl); err != nil {
		return nil, nil, fmt.Errorf("failed to update time log: %w", err)
	}

	updatedLog, err := model.GetTimeLogByID(server.db, args.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve updated time log: %w", err)
	}

	response := map[string]interface{}{
		"id":          *updatedLog.ID,
		"category_id": updatedLog.CategoryID,
		"start_time":  formatSGDateTime(updatedLog.StartTime),
		"remark":      updatedLog.Remark,
	}
	if updatedLog.EndTime != nil {
		response["end_time"] = formatSGDateTime(*updatedLog.EndTime)
	}
	if updatedLog.TaskID != nil {
		response["task_id"] = *updatedLog.TaskID
	}

	return formatMCPResponse("Time log updated successfully", response)
}
