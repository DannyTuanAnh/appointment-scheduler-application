package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type DealershipRoutes struct {
	dealership_handler *handlers.DealershipHandler
}

func NewDealershipRoutes(handler *handlers.DealershipHandler) Routes {
	return &DealershipRoutes{
		dealership_handler: handler,
	}
}

func (ar *DealershipRoutes) RegisterApp(r *gin.RouterGroup) {
	dealership := r.Group("/dealership")
	{
		dealership.GET("", ar.dealership_handler.GetAllDealerships)
		dealership.GET("/search", ar.dealership_handler.SearchDealerships)
		dealership.GET("/:id", ar.dealership_handler.GetDealershipByID)
	}
}

func (ar *DealershipRoutes) RegisterAdmin(r *gin.RouterGroup) {
	dealership := r.Group("/dealership")
	{
		dealership.POST("", ar.dealership_handler.CreateDealership)
		dealership.PUT("/:id", ar.dealership_handler.UpdateDealership)
		dealership.DELETE("/:id", ar.dealership_handler.DeleteDealership)
	}
}
