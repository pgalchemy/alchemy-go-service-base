package instrumentation

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve serves the metrics endpoint on the provided router
func Serve(e *gin.Engine) {
	// Use metrics gathering middleware
	e.Use(Middleware)

	// Attach metrics handler
	ph := promhttp.Handler()
	e.Handle("GET", "/metrics", func(c *gin.Context) {
		ph.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})
}
