package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

// ListGroundTask GET /api/ground-task
func ListGroundTask(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListTasks())
}
