package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Book-specific counter
	BookRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "book_requests_total",
			Help: "Total number of requests to /books endpoint",
		},
		[]string{"method"},
	)

	// Global HTTP requests counter
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Histogram for request duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // default: [0.005, 0.01, ... 10]
		},
		[]string{"method", "path"},
	)
)

// InitMetrics registers all metrics
func InitMetrics() {
	prometheus.MustRegister(BookRequests, HTTPRequestsTotal, HTTPRequestDuration)
}

// MetricsHandler returns a Gin handler for /metrics
func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// Helper: normalize status code (e.g., 200, 404, 500)
func httpStatusToLabel(code int) string {
	return fmt.Sprintf("%d", code)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()

		// Update global metrics
		HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(), // registered path, e.g. "/books"
			httpStatusToLabel(statusCode),
		).Inc()

		HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
