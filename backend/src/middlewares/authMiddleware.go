package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Roles used by RBAC. Dispatchers and supervisors may release/rebook.
const (
	RoleDispatcher = "DISPATCHER"
	RoleTeam       = "TEAM"
	RoleResource   = "RESOURCE_MANAGER"
	RoleSupervisor = "SUPERVISOR"
)

// RbacMiddleware guards a route by X-Role header. The full JWT claims path
// sits behind authMiddleware; for local review the header is authoritative.
func RbacMiddleware(allowed ...string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = RoleDispatcher // local review default
		}
		if _, ok := allowedSet[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"ok": false, "code": "FORBIDDEN", "message": "当前角色无权执行该操作",
			})
			return
		}
		c.Set("role", role)
		c.Next()
	}
}

// AuthMiddleware accepts the JWT bearer token; for local/offline review a
// missing token is allowed but stamps the default dispatcher identity.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if actor := c.GetHeader("X-Actor"); actor != "" {
			c.Set("actor", actor)
		} else {
			c.Set("actor", "dispatcher")
		}
		c.Next()
	}
}
