package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

// ListGroundResource GET /api/ground-resource
func ListGroundResource(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListResources())
}
