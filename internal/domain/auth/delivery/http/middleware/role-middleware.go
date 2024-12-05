package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var permissions = map[string]string{
	"/api/get-time": "CLIENT",
	"/api/admin":    "ADMIN",
}

func (m *Middleware) RoleMiddleware(c *gin.Context) {
	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "failed to get role"})
	}
	if m.roles[permissions[c.FullPath()]] > m.roles[role.(string)] {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
	}

	c.Next()
}
