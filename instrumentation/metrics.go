package instrumentation

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Predefined set of metrics to export
var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "How many HTTP requests processed, partitioned by status code, HTTP method, and path.",
		},
		[]string{"code", "method", "path"},
	)

	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "request_duration_seconds",
			Help: "The HTTP request latencies in seconds, partitioned by status code, HTTP method, and path.",
		},
		[]string{"code", "method", "path"},
	)

	requestSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "request_size_bytes",
			Help: "The HTTP request sizes in bytes, partitioned by status code, HTTP method, and path.",
		},
		[]string{"code", "method", "path"},
	)

	responseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "response_size_bytes",
			Help: "The HTTP response sizes in bytes, partitioned by status code, HTTP method, and path.",
		},
		[]string{"code", "method", "path"},
	)
)

func init() {
	// Register all metrics with prometheus
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestSize)
	prometheus.MustRegister(responseSize)
}

// Middleware provides a gin middleware for collecting metric data
func Middleware(c *gin.Context) {
	// dont track metrics requests
	if c.Request.URL.Path == "/metrics" {
		c.Next()
		return
	}

	// Mark start time
	start := time.Now()

	// Asynchronously fetch request size
	reqSize := make(chan int)
	go getRequestSize(c.Request, reqSize)

	// Wait for request to complete
	c.Next()

	status := strconv.Itoa(c.Writer.Status())
	elapsed := float64(time.Since(start)) / float64(time.Second)
	resSize := float64(c.Writer.Size())

	// Update metrics
	requestCount.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Inc()
	requestDuration.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
	requestSize.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(float64(<-reqSize))
	responseSize.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(resSize)
}

// Get total request size
func getRequestSize(r *http.Request, out chan int) {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	out <- s
}
