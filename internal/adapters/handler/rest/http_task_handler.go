package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/jittakal/todo-api-hexagonal/api/rest"
	"github.com/jittakal/todo-api-hexagonal/internal/domain"
	"github.com/jittakal/todo-api-hexagonal/internal/errors"
	"github.com/jittakal/todo-api-hexagonal/internal/middleware"
)

var (
	logger = middleware.NewLogger()
)

type GinTaskHandler struct {
	repo domain.TaskRepository
}

var _ api.TaskHandler = (*GinTaskHandler)(nil)

func NewGinTaskHandler(repository domain.TaskRepository) *GinTaskHandler {
	return &GinTaskHandler{
		repo: repository,
	}
}

// CreateTask handles the creation of a new task.
// @Summary Create a new task
// @Description Create a new task
// @Accept json
// @Produce json
// @Param request body TaskCreateRequest true "Task creation request"
// @Success 201 {object} TaskCreateResponse
// @Router /todo/v1/tasks [post]
func (h *GinTaskHandler) CreateTask(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	logger.Info(ctx, "Request received to create new task")

	var newTask *api.TaskCreateRequest
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	taskCreateRes, err := h.repo.Create(ctx, newTask.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, taskCreateRes)
}

// GetTaskHandler gets details of a specific task.
// @Summary Get details of a task
// @Description Get details of a task by ID
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
// @Failure 404 "Task Not Found"
// @Router /todo/v1/tasks/{id} [get]
func (h *GinTaskHandler) GetTask(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	taskID := c.Param("id")

	task, err := h.repo.GetByID(ctx, taskID)
	if err != nil {
		if err == errors.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTaskHandler updates details of a specific task.
// @Summary Update details of a task
// @Description Update details of a task by ID
// @Accept json
// @Param id path string true "Task ID"
// @Param request body domain.Task true "Task update request"
// @Success 200 "OK"
// @Failure 404 "Task Not Found"
// @Router /todo/v1/tasks/{id} [put]
func (h *GinTaskHandler) UpdateTask(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	taskID := c.Param("id")

	var updatedTask *domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedTask.ID = taskID
	if err := h.repo.Update(ctx, updatedTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteTask deletes a specific task.
// @Summary Delete a task
// @Description Delete a task by ID
// @Param id path string true "Task ID"
// @Success 200 "OK"
// @Failure 404 "Task Not Found"
// @Router /todo/v1/tasks/{id} [delete]
func (h *GinTaskHandler) DeleteTask(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	taskID := c.Param("id")

	if err := h.repo.Delete(ctx, taskID); err != nil {
		if err == errors.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetAllTasks gets a list of all tasks.
// @Summary Get all tasks
// @Description Get a list of all tasks
// @Produce json
// @Success 200 {array} domain.Task
// @Failure 404 "Task Not Found"
// @Router /todo/v1/tasks [get]
func (h *GinTaskHandler) GetAllTasks(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	logger.Info(ctx, "Request received to get all tasks ...")

	tasks, err := h.repo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// MarkTaskDone marks a task as done.
// @Summary Mark task as done
// @Description Mark a task as done by ID
// @Param id path string true "Task ID"
// @Success 200 "OK"
// @Failure 404 "Task Not Found"
// @Router /todo/v1/tasks/{id}/done [post]
func (h *GinTaskHandler) MarkTaskDone(c *gin.Context) {
	logID := middleware.GenerateLogID()
	ctx := middleware.WithCorrelationLogID(c, logID)

	taskID := c.Param("id")

	task, err := h.repo.GetByID(ctx, taskID)

	if err == errors.ErrTaskNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Done = true
	h.repo.Update(ctx, task)
	c.Status(http.StatusOK)
}
