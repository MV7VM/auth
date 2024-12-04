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
	if permissions[c.FullPath()] != role.(string) {
		c.AbortWithStatus(http.StatusForbidden)
	}

	c.Next()
}
