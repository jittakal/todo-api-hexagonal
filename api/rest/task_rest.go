package rest

import "github.com/gin-gonic/gin"

type TaskHandler interface {
	CreateTask(c *gin.Context)

	GetTask(c *gin.Context)

	UpdateTask(c *gin.Context)

	DeleteTask(c *gin.Context)

	GetAllTasks(c *gin.Context)

	MarkTaskDone(c *gin.Context)
}
