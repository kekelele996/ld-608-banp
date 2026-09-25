package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

// fail maps a service BusinessError onto a structured HTTP response.
// Rejected releases return 409 with the full violation list so the UI can
// highlight exact tasks/delays/bookings; other business errors return 400/404.
func fail(c *gin.Context, err error) {
	var bizErr *services.BusinessError
	if errors.As(err, &bizErr) {
		status := http.StatusBadRequest
		switch bizErr.Code {
		case "TURNAROUND_NOT_FOUND", "TASK_NOT_FOUND", "BOOKING_NOT_FOUND", "RESOURCE_NOT_FOUND":
			status = http.StatusNotFound
		case "RELEASE_REJECTED", "TURNAROUND_NOT_RELEASABLE":
			status = http.StatusConflict
		}
		body := gin.H{
			"ok":      false,
			"code":    bizErr.Code,
			"message": bizErr.Message,
		}
		if len(bizErr.Violations) > 0 {
			body["violations"] = bizErr.Violations
		}
		if bizErr.Violation != nil {
			body["violation"] = bizErr.Violation
		}
		c.AbortWithStatusJSON(status, body)
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"ok": false, "code": "INTERNAL_ERROR", "message": err.Error(),
	})
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": data})
}

// actor resolves the acting user from X-Actor header or JSON body fallback.
// Full JWT/RBAC lives behind authMiddleware; this keeps local review working.
func actor(c *gin.Context) string {
	if a := c.GetHeader("X-Actor"); a != "" {
		return a
	}
	return "dispatcher"
}
