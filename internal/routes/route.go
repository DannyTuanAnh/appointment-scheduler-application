package routes

import (
	"context"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(ctx context.Context, r *gin.Engine, db sqlc.Querier, routes ...Routes) {
	// Register middleware for all routes
	r.Use(middlewares.LoggerMiddleware())

	// Public health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	base := r.Group("/api/v1")

	public := base.Group("/app")
	for _, route := range routes {
		if publicRoute, ok := route.(interface{ RegisterPublic(r *gin.RouterGroup) }); ok {
			publicRoute.RegisterPublic(public)
		}
	}

	admin := base.Group("/admin")
	for _, route := range routes {
		if adminRoute, ok := route.(interface{ RegisterAdmin(r *gin.RouterGroup) }); ok {
			adminRoute.RegisterAdmin(admin)
		}
	}
}
