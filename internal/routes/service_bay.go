package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type ServiceBayRoutes struct {
	handler *handlers.ServiceBayHandler
}

func NewServiceBayRoutes(handler *handlers.ServiceBayHandler) Routes {
	return &ServiceBayRoutes{handler: handler}
}

func (sr *ServiceBayRoutes) RegisterApp(r *gin.RouterGroup) {
	serviceBay := r.Group("/service-bay")
	{
		serviceBay.GET("", sr.handler.ListServiceBays)
		serviceBay.GET("/search", sr.handler.SearchServiceBays)
		serviceBay.GET("/:id", sr.handler.GetServiceBayByID)

		serviceBayType := serviceBay.Group("/type")
		{
			serviceBayType.GET("", sr.handler.ListServiceBayTypes)
			serviceBayType.GET("/search", sr.handler.SearchServiceBayTypes)
			serviceBayType.GET("/:id", sr.handler.GetServiceBayTypeByID)
		}
	}
}

func (sr *ServiceBayRoutes) RegisterAdmin(r *gin.RouterGroup) {
	serviceBay := r.Group("/service-bay")
	{
		serviceBay.POST("", sr.handler.CreateServiceBay)
		serviceBay.PUT("/:id", sr.handler.UpdateServiceBay)
		serviceBay.DELETE("/:id", sr.handler.DeleteServiceBay)

		serviceBayType := serviceBay.Group("/type")
		{
			serviceBayType.POST("", sr.handler.CreateServiceBayType)
			serviceBayType.PUT("/:id", sr.handler.UpdateServiceBayType)
			serviceBayType.DELETE("/:id", sr.handler.DeleteServiceBayType)
		}
	}
}
