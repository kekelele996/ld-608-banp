package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLogMiddleware writes an access line per mutating request. Domain-level
// audit entries (release/rebook) are written by services inside transactions;
// this middleware covers the transport-level trail.
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Request.Method == http.MethodGet {
			return
		}
		actor := c.GetHeader("X-Actor")
		if actor == "" {
			actor = "anonymous"
		}
		log.Printf("[audit] actor=%s method=%s path=%s status=%d dur=%s",
			actor, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}
