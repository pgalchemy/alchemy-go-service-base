package logging

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pgalchemy/alchemy-go-service-base/scope"
	"github.com/sirupsen/logrus"
)

type (
	// AccessLogMiddlewareOptions configures the AccessLogMiddleware
	AccessLogMiddlewareOptions struct {
		Ignore []string
	}
)

// New provides a new logger instance
func New(f logrus.Formatter) *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.Level = logrus.DebugLevel
	l.Formatter = f
	return l
}

// ErrorLogMiddleware provides a gin middleware for logging errors
func ErrorLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			rs := scope.GetRequestScope(c)
			rs.Logger.WithField("error", c.Errors[0]).Warn("error: ", c.Errors[0].Error())
		}
	}
}

// AccessLogMiddleware provides a gin middleware for logging requests
func AccessLogMiddleware(opts *AccessLogMiddlewareOptions) gin.HandlerFunc {
	if opts.Ignore == nil {
		opts.Ignore = []string{}
	}

	return func(c *gin.Context) {
		if opts.ignored(c.Request.URL.Path) {
			return
		}
		rs := scope.GetRequestScope(c)

		// Get start time
		start := time.Now()
		path := c.Request.URL.Path

		// Continue
		c.Next()

		if path == "/metrics" {
			return
		}

		end := time.Now()
		latency := end.Sub(start) / time.Millisecond
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		httpVersion := fmt.Sprintf("%d.%d", c.Request.ProtoMajor, c.Request.ProtoMinor)

		rs.Logger.WithFields(logrus.Fields{
			"remote-address": clientIP,
			"ip":             clientIP,
			"method":         method,
			"path":           path,
			"url":            c.Request.Host + path,
			"referrer":       c.Request.Header.Get("referrer"),
			"user-agent":     c.Request.Header.Get("user-agent"),
			"http-version":   httpVersion,
			"response-time":  latency,
			"response-size":  c.Writer.Size(),
			"status-code":    statusCode,
			"req-headers":    formatHeaders(c.Request.Header),
			"res-headers":    formatHeaders(c.Writer.Header()),
		}).Info(fmt.Sprintf("%s %s %d", method, path, statusCode))
	}
}

func (opts *AccessLogMiddlewareOptions) ignored(path string) bool {
	if opts.Ignore == nil || len(opts.Ignore) == 0 {
		return false
	}

	for _, v := range opts.Ignore {
		if path == v {
			return true
		}
	}

	return false
}

func formatHeaders(h http.Header) map[string]string {
	formatted := map[string]string{}
	for k, v := range h {
		if len(v) > 0 {
			formatted[k] = v[0]
		}
	}
	return formatted
}
