package middlewares

import (
	"strconv"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/metrics"
	"github.com/gin-gonic/gin"
)

type metricResponseWriter struct {
	gin.ResponseWriter
	size int
}

func (w *metricResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		writer := &metricResponseWriter{
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = writer

		ctx.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ctx.Writer.Status())

		path := ctx.FullPath()
		if path == "" {
			path = ctx.Request.URL.Path
		}

		method := ctx.Request.Method

		metrics.HTTPRequestTotal.WithLabelValues(method, path, status).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		metrics.HTTPRequestSizeBytes.WithLabelValues(method, path).
			Observe(float64(ctx.Request.ContentLength))

		metrics.HTTPResponseSizeBytes.WithLabelValues(method, path, status).
			Observe(float64(writer.size))

		if ctx.Writer.Status() >= 400 {
			metrics.HTTPErrorTotal.WithLabelValues(method, path, status).Inc()
		}
	}
}
