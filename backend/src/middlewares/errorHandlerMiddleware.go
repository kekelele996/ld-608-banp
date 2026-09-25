package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware recovers from panics and emits the same error
// envelope used by controllers.fail so the frontend never sees a bare 500.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v %s", rec, c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"ok": false, "code": "INTERNAL_ERROR", "message": "服务内部错误",
				})
			}
		}()
		c.Next()
	}
}
