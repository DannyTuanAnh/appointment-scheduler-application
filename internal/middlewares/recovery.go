package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/metrics"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/logs"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	logger := logs.InitLogger(utils.GetEnv("PATH_LOGGER_PANIC", ""))

	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Increase panic metric
				metrics.HTTPPanicsTotal.Inc()

				// Log panic + stack trace
				logger.Error().
					Str("method", ctx.Request.Method).
					Str("path", ctx.Request.URL.Path).
					Str("client_ip", ctx.ClientIP()).
					Interface("panic", err).
					Bytes("stack_trace", debug.Stack()).
					Msg("Panic recovered")

				// Prevent duplicated write
				if !ctx.Writer.Written() {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "Internal Server Error",
						"error":   fmt.Sprintf("%v", err),
					})
				} else {
					ctx.Abort()
				}
			}
		}()

		ctx.Next()
	}
}
