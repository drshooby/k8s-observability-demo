package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	tasks      = make(map[int]Task)
	nextID     = 1
	tasksMutex sync.Mutex
)

func main() {
	router := gin.Default()

	// Health and readiness endpoints
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "READY")
	})

	// Middleware to simulate delay and error
	router.Use(func(c *gin.Context) {
		if delayStr := c.Query("delay"); delayStr != "" {
			if delayMs, err := strconv.Atoi(delayStr); err == nil && delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}
		}
		if errStr := c.Query("error"); errStr == "true" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "simulated error"})
			return
		}
		c.Next()
	})

	// List all tasks
	router.GET("/tasks", func(c *gin.Context) {
		tasksMutex.Lock()
		defer tasksMutex.Unlock()

		list := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			list = append(list, t)
		}
		c.JSON(http.StatusOK, list)
	})

	// Add a new task (expects JSON with {"name": "task name"})
	router.POST("/tasks", func(c *gin.Context) {
		var newTask struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		tasksMutex.Lock()
		task := Task{ID: nextID, Name: newTask.Name}
		tasks[nextID] = task
		nextID++
		tasksMutex.Unlock()

		c.JSON(http.StatusCreated, task)
	})

	// Delete a task by ID
	router.DELETE("/tasks/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		tasksMutex.Lock()
		defer tasksMutex.Unlock()

		if _, exists := tasks[id]; !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		delete(tasks, id)
		c.Status(http.StatusNoContent)
	})

	router.Run(":8080")
}
