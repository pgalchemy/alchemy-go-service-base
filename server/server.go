package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pgsuite/alchemy-go-service-base/errors"
	"github.com/pgsuite/alchemy-go-service-base/logging"
	"github.com/pgsuite/alchemy-go-service-base/scope"
	"github.com/sirupsen/logrus"
)

type (
	// Config represents the server configuratino
	Config struct {
		// Logger is the logrus logging instance to utilize in middlewares
		Logger *logrus.Logger

		// IgnoredPaths omits access logs for provided paths
		IgnoredPaths []string
	}
)

// New provides a new gin server instance
func New(c *Config) *gin.Engine {
	// Create server instance
	e := gin.New()

	ignored := []string{"/", ""}
	if c.IgnoredPaths != nil {
		ignored = c.IgnoredPaths
	}

	// Attach middlewares
	if c.Logger != nil {
		e.Use(scope.Middleware(c.Logger))
		e.Use(logging.AccessLogMiddleware(&logging.AccessLogMiddlewareOptions{Ignore: ignored}))
		e.Use(logging.ErrorLogMiddleware())
	}

	e.Use(errors.Middleware)

	return e
}
