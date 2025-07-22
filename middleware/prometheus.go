package mw

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
Prometheus Middleware to record the following metrics:
- Duration of the requests(with: code, handler, method).
- Count of the requests(with: code, handler, method).
- Size of the responses(with: code, handler, method).
- Number requests being handled concurrently at a given time a.k.a inflight requests (with: handler).
*/

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "handler", "code"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler", "code"},
	)

	inflightGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Current number of inflight requests.",
		},
		[]string{"handler"},
	)

	responseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_response_size_bytes",
			Help:       "Size of HTTP responses.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "handler", "code"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, inflightGauge, responseSize)
}

// PrometheusMiddleware returns a Gin middleware for metrics instrumentation
func PrometheusMiddleware(handlerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method

		inflightGauge.WithLabelValues(handlerName).Inc()
		defer inflightGauge.WithLabelValues(handlerName).Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()
		size := float64(c.Writer.Size())

		labels := prometheus.Labels{
			"method":  method,
			"handler": handlerName,
			"code":    fmt.Sprintf("%d", statusCode),
		}

		requestCounter.With(labels).Inc()
		requestDuration.With(labels).Observe(duration)
		responseSize.With(labels).Observe(size)
	}
}

func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
