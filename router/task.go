package router

import (
	"net/http"
	"strconv"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/gin-gonic/gin"
)

// 添加任务相关路由
func setupTaskRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/tasks", listTasksHandler(deps))
	group.POST("/tasks", createTaskHandler(deps))
	group.GET("/tasks/:id", getTaskHandler(deps))
	group.PUT("/tasks/:id", updateTaskHandler(deps))
	group.DELETE("/tasks/:id", deleteTaskHandler(deps))
	group.POST("/tasks/:id/complete", completeTaskHandler(deps))
	group.POST("/tasks/:id/incomplete", incompleteTaskHandler(deps))
	group.POST("/tasks/:id/suspend", suspendTaskHandler(deps))
	group.POST("/tasks/:id/unsuspend", unsuspendTaskHandler(deps))
	group.GET("/tasks/stats/:date", getTaskStatsHandler(deps))
}

func createTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if err := deps.Service.CreateTask(task); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整的Tag信息
		if createdTask, err := deps.Service.GetTaskByID(*task.ID); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(createdTask), "Task created successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(task), "Task created successfully"))
		}
	}
}

func listTasksHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
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
				tasks, err = deps.Service.GetTasksByDate(date, includeSuspended, includeCompleted)
			}
		} else {
			tasks, err = deps.Service.GetAllTasks(includeSuspended, includeCompleted)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.TasksToProto(tasks), "Tasks retrieved successfully"))
	}
}

func getTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		task, err := deps.Service.GetTaskByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(task), "Task retrieved successfully"))
	}
}

func updateTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		// 先检查任务是否存在
		existingTask, err := deps.Service.GetTaskByID(id)
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

		if err := deps.Service.UpdateTask(existingTask); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		if updatedTask, err := deps.Service.GetTaskByID(id); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(updatedTask), "Task updated successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.TaskToProto(existingTask), "Task updated successfully"))
		}
	}
}

func deleteTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		// 先检查任务是否存在
		if _, err := deps.Service.GetTaskByID(id); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
			return
		}

		if err := deps.Service.DeleteTask(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Task deleted successfully"))
	}
}

func completeTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		if err := deps.Service.MarkTaskAsCompleted(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as completed"))
	}
}

func incompleteTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		if err := deps.Service.MarkTaskAsIncomplete(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Task marked as incomplete"))
	}
}

func suspendTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		// 先检查任务是否存在
		if _, err := deps.Service.GetTaskByID(id); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
			return
		}

		if err := deps.Service.SuspendTask(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Task suspended successfully"))
	}
}

func unsuspendTaskHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid task ID"))
			return
		}
		id := int32(id64)

		// 先检查任务是否存在
		if _, err := deps.Service.GetTaskByID(id); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Task not found"))
			return
		}

		if err := deps.Service.UnsuspendTask(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Task unsuspended successfully"))
	}
}

func getTaskStatsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD"))
			return
		}

		stats, err := deps.Service.GetTaskStats(date)
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
}
