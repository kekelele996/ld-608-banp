package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

// ListDelayEvent GET /api/delay-event
func ListDelayEvent(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListDelays())
}
