package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	router := gin.Default()

	// Health and readiness endpoints
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "READY")
	})

	// Summary endpoint: fetch tasks from service1 and return count
	router.GET("/summary", func(c *gin.Context) {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("http://service1:8080/tasks")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to get tasks from service1"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service1 returned non-200 status"})
			return
		}

		var tasks []Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode tasks"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"task_count": len(tasks),
		})
	})

	router.Run(":8081")
}
