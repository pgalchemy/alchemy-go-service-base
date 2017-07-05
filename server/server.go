package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pgalchemy/alchemy-go-service-base/logging"
	"github.com/pgalchemy/alchemy-go-service-base/scope"
	"github.com/sirupsen/logrus"
)

type (
	// Config represents the server configuratino
	Config struct {
		// Logger is the logrus logging instance to utilize in middlewares
		Logger *logrus.Logger
	}
)

// New provides a new gin server instance
func New(c *Config) *gin.Engine {
	// Create server instance
	e := gin.New()

	// Attach middlewares
	e.Use(scope.Middleware(c.Logger))
	e.Use(logging.AccessLogMiddleware())

	return e
}
