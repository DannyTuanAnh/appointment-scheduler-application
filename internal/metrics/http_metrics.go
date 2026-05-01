package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var once sync.Once

var (
	// Total requests
	HTTPRequestTotal *prometheus.CounterVec

	// Request duration (latency)
	HTTPRequestDuration *prometheus.HistogramVec

	// Requests currently processing
	HTTPRequestsInFlight prometheus.Gauge

	// Response size
	HTTPResponseSizeBytes *prometheus.HistogramVec

	// Request size
	HTTPRequestSizeBytes *prometheus.HistogramVec

	// Panic count
	HTTPPanicsTotal prometheus.Counter

	// Error requests
	HTTPErrorTotal *prometheus.CounterVec
)

func Init() {
	once.Do(func() {
		HTTPRequestTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)

		HTTPRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		)

		HTTPRequestsInFlight = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of in-flight requests",
			},
		)

		HTTPRequestSizeBytes = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "Size of HTTP request in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 2, 10),
			},
			[]string{"method", "path"},
		)

		HTTPResponseSizeBytes = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "Size of HTTP response in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 2, 10),
			},
			[]string{"method", "path", "status"},
		)

		HTTPPanicsTotal = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_panics_total",
				Help: "Total number of recovered panics",
			},
		)

		HTTPErrorTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of error responses",
			},
			[]string{"method", "path", "status"},
		)

		prometheus.MustRegister(
			HTTPRequestTotal,
			HTTPRequestDuration,
			HTTPRequestsInFlight,
			HTTPRequestSizeBytes,
			HTTPResponseSizeBytes,
			HTTPPanicsTotal,
			HTTPErrorTotal,
		)
	})
}
