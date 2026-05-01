package routes

import (
	"context"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Routes interface {
	RegisterApp(r *gin.RouterGroup)
}

func RegisterRoutes(ctx context.Context, r *gin.Engine, routes ...Routes) {
	// Register middleware for all routes
	r.Use(middlewares.RecoveryMiddleware(), middlewares.MetricsMiddleware(), middlewares.LoggerMiddleware())

	// Public health check endpoint
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	// Public metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	base := r.Group("/api/v1")

	public := base.Group("/app")
	for _, route := range routes {
		if publicRoute, ok := route.(interface{ RegisterApp(r *gin.RouterGroup) }); ok {
			publicRoute.RegisterApp(public)
		}
	}

	admin := base.Group("/admin")
	for _, route := range routes {
		if adminRoute, ok := route.(interface{ RegisterAdmin(r *gin.RouterGroup) }); ok {
			adminRoute.RegisterAdmin(admin)
		}
	}
}
