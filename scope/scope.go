package scope

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	// RequestScope contains request specific utilities
	RequestScope struct {
		Logger *logrus.Entry
		Config map[string]string
	}
)

// Middleware provides a gin middleware for attaching request scope
func Middleware(l *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to read in request id from gateway
		rid := c.Request.Header.Get("X-Request-Id")

		// Create dedicated request logger
		requestLogger := l.WithFields(logrus.Fields{"request_id": rid})

		// Create and attach request scope
		rs := &RequestScope{Logger: requestLogger}
		c.Set("scope", rs)
		c.Next()
	}
}

// GetRequestScope provides a utility for retrieving request scope
// out of context.
func GetRequestScope(c *gin.Context) *RequestScope {
	scope, _ := c.Get("scope")
	return scope.(*RequestScope)
}
