package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var permissions = map[string]string{
	"/api/get-time": "CLIENT",
	"/api/admin":    "ADMIN",
}

var roles = map[string]int{
	"CLIENT": 0,
	"ADMIN":  1,
}

func (m *Middleware) RoleMiddleware(c *gin.Context) {
	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "failed to get role"})
	}
	if roles[permissions[c.FullPath()]] > roles[role.(string)] {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
	}

	c.Next()
}
