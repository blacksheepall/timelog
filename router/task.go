package router

import (
	"net/http"
	"strconv"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// 添加任务相关路由
func setupTaskRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/tasks", listTasksHandler)
	group.POST("/tasks", createTaskHandler)
	group.GET("/tasks/:id", getTaskHandler)
	group.PUT("/tasks/:id", updateTaskHandler)
	group.DELETE("/tasks/:id", deleteTaskHandler)
	group.POST("/tasks/:id/complete", completeTaskHandler)
	group.POST("/tasks/:id/incomplete", incompleteTaskHandler)
	group.POST("/tasks/:id/suspend", suspendTaskHandler)
	group.POST("/tasks/:id/unsuspend", unsuspendTaskHandler)
	group.GET("/tasks/stats/:date", getTaskStatsHandler)
}

// CreateTaskHandler godoc
// @Summary 创建任务
// @Description 新增一项任务
// @Tags task
// @Accept json
// @Produce json
// @Param data body gen.Task true "任务数据"
// @Success 200 {object} gen.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks [post]
func createTaskHandler(c *gin.Context) {
	var req timelogv1.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	task, err := mapper.TaskFromCreateRequest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := service.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整的Tag信息
	if createdTask, err := service.GetTaskByID(*task.ID); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(createdTask), "Task created successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(task), "Task created successfully"))
	}
}

// ListTasksHandler godoc
// @Summary 获取任务列表
// @Description 获取所有任务，支持按日期过滤、是否包含暂停的任务和是否包含已完成的任务
// @Tags task
// @Produce json
// @Param date query string false "日期过滤 (YYYY-MM-DD格式)"
// @Param include_suspended query boolean false "是否包含暂停的任务 (默认false)"
// @Param include_completed query boolean false "是否包含已完成的任务 (默认false)"
// @Success 200 {array} gen.Task
// @Failure 500 {object} map[string]string
// @Router /api/tasks [get]
func listTasksHandler(c *gin.Context) {
	dateStr := c.Query("date")
	includeSuspended := c.Query("include_suspended") == "true"
	includeCompleted := c.Query("include_completed") == "true"

	var tasks []gen.Task
	var err error

	if dateStr != "" {
		// 解析日期
		if date, parseErr := time.Parse("2006-01-02", dateStr); parseErr != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD"))
			return
		} else {
			tasks, err = service.GetTasksByDate(date, includeSuspended, includeCompleted)
		}
	} else {
		tasks, err = service.GetAllTasks(includeSuspended, includeCompleted)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(mapper.TasksToProto(tasks), "Tasks retrieved successfully"))
}

// GetTaskHandler godoc
// @Summary 获取单个任务
// @Description 根据ID获取任务详情
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} gen.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [get]
func getTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	task, err := service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(task), "Task retrieved successfully"))
}

// UpdateTaskHandler godoc
// @Summary 更新任务
// @Description 更新任务信息
// @Tags task
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param data body gen.Task true "任务数据"
// @Success 200 {object} gen.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [put]
func updateTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	// 先检查任务是否存在
	existingTask, err := service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}

	var req timelogv1.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := mapper.ApplyTaskUpdate(existingTask, &req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := service.UpdateTask(existingTask); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整信息
	if updatedTask, err := service.GetTaskByID(id); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(updatedTask), "Task updated successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(existingTask), "Task updated successfully"))
	}
}

// DeleteTaskHandler godoc
// @Summary 删除任务
// @Description 删除指定任务
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id} [delete]
func deleteTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	// 先检查任务是否存在
	if _, err := service.GetTaskByID(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}

	if err := service.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Task deleted successfully"))
}

// CompleteTaskHandler godoc
// @Summary 标记任务为完成
// @Description 将任务标记为完成状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/complete [post]
func completeTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	if err := service.MarkTaskAsCompleted(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as completed"))
}

// IncompleteTaskHandler godoc
// @Summary 标记任务为未完成
// @Description 将任务标记为未完成状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/incomplete [post]
func incompleteTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	if err := service.MarkTaskAsIncomplete(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as incomplete"))
}

// SuspendTaskHandler godoc
// @Summary 暂停任务
// @Description 将任务标记为暂停状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/suspend [post]
func suspendTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	// 先检查任务是否存在
	if _, err := service.GetTaskByID(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}

	if err := service.SuspendTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Task suspended successfully"))
}

// UnsuspendTaskHandler godoc
// @Summary 取消暂停任务
// @Description 将任务取消暂停状态
// @Tags task
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/unsuspend [post]
func unsuspendTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
		return
	}
	id := int32(id64)

	// 先检查任务是否存在
	if _, err := service.GetTaskByID(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
		return
	}

	if err := service.UnsuspendTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Task unsuspended successfully"))
}

// GetTaskStatsHandler godoc
// @Summary 获取任务统计
// @Description 获取指定日期的任务完成统计
// @Tags task
// @Produce json
// @Param date path string true "日期 (YYYY-MM-DD格式)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/stats/{date} [get]
func getTaskStatsHandler(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD"))
		return
	}

	stats, err := service.GetTaskStats(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	payload, err := mapper.TaskStatsToProto(stats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(payload, "Task stats retrieved successfully"))
}
