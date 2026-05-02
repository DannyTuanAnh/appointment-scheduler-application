package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type AppointmentRoutes struct {
	appointment_handler *handlers.AppointmentHandler
}

func NewAppointmentRoutes(handler *handlers.AppointmentHandler) Routes {
	return &AppointmentRoutes{
		appointment_handler: handler,
	}
}

func (ar *AppointmentRoutes) RegisterApp(r *gin.RouterGroup) {
	appointment := r.Group("/appointment")
	{
		appointment.GET("", ar.appointment_handler.GetAppointment)
	}
}

func (ar *AppointmentRoutes) RegisterAdmin(r *gin.RouterGroup) {
}
