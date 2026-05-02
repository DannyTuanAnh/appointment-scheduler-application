package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type ServiceRoutes struct {
	handler *handlers.ServiceHandler
}

func NewServiceRoutes(handler *handlers.ServiceHandler) Routes {
	return &ServiceRoutes{handler: handler}
}

func (sr *ServiceRoutes) RegisterApp(r *gin.RouterGroup) {
	service := r.Group("/service")
	{
		service.GET("", sr.handler.ListServices)
		service.GET("/search", sr.handler.SearchServices)
		service.GET("/:id", sr.handler.GetServiceDetailByID)
	}
}

func (sr *ServiceRoutes) RegisterAdmin(r *gin.RouterGroup) {
	service := r.Group("/service")
	{
		service.POST("", sr.handler.CreateService)
		service.PUT("/:id", sr.handler.UpdateService)
		service.DELETE("/:id", sr.handler.DeleteService)

		service.POST("/:id/requirements", sr.handler.AddSkillRequirementsToService)
		service.DELETE("/:id/requirements", sr.handler.RemoveSkillRequirementsFromService)
	}
}
