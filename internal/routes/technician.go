package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type TechnicianRoutes struct {
	handler *handlers.TechnicianHandler
}

func NewTechnicianRoutes(handler *handlers.TechnicianHandler) Routes {
	return &TechnicianRoutes{handler: handler}
}

func (tr *TechnicianRoutes) RegisterApp(r *gin.RouterGroup) {
	tech := r.Group("/technician")
	{
		tech.GET("", tr.handler.ListTechnicians)
		tech.GET("/search", tr.handler.SearchTechnicians)
		tech.GET("/active", tr.handler.FindActiveTechnicians)
		tech.GET("/:id", tr.handler.GetTechnicianByID)
		tech.GET("/:id/detail", tr.handler.GetDetailTechnicianByID)
	}
}

func (tr *TechnicianRoutes) RegisterAdmin(r *gin.RouterGroup) {
	tech := r.Group("/technician")
	{
		tech.POST("", tr.handler.CreateTechnician)
		tech.PUT("/:id", tr.handler.UpdateTechnician)
		tech.DELETE("/:id", tr.handler.DeleteTechnician)

		tech.PUT("/:id/leave", tr.handler.SetTechnicianOnLeave)
		tech.PUT("/:id/back", tr.handler.SetTechnicianBackToWork)
		tech.PUT("/:id/transfer", tr.handler.TransferTechnicianDealership)

		tech.POST("/:id/skills", tr.handler.AddSkillsToTechnician)
		tech.DELETE("/:id/skills", tr.handler.RemoveSkillsFromTechnician)
	}
}
