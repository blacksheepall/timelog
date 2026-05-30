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
