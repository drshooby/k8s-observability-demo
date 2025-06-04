package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	tasks      = make(map[int]Task)
	nextID     = 1
	tasksMutex sync.Mutex

	// Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests processed, labeled by path and method.",
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
}

func main() {
	router := gin.Default()

	// Register /metrics for Prometheus scraping
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Middleware: record metrics and simulate delay/error
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path // fallback
		}
		method := c.Request.Method

		// Simulated delay
		if delayStr := c.Query("delay"); delayStr != "" {
			if delayMs, err := strconv.Atoi(delayStr); err == nil && delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}
		}

		// Simulated error
		if c.Query("error") == "true" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "simulated error"})
			return
		}

		c.Next()

		// Record request after response
		httpRequestsTotal.WithLabelValues(path, method).Inc()
		_ = time.Since(start) // useful later for latency histograms
	})

	// Health checks
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "READY")
	})

	// Task handlers
	router.GET("/tasks", func(c *gin.Context) {
		tasksMutex.Lock()
		defer tasksMutex.Unlock()
		list := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			list = append(list, t)
		}
		c.JSON(http.StatusOK, list)
	})

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
