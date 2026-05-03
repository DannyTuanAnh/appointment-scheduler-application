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
		// Availability (existing)
		appointment.GET("", ar.appointment_handler.GetAppointment)

		// Create + status updates
		appointment.POST("", ar.appointment_handler.CreateAppointment)
		appointment.PUT("/:id/cancel", ar.appointment_handler.CancelledAppointment)

	}
}

func (ar *AppointmentRoutes) RegisterAdmin(r *gin.RouterGroup) {
	appointment := r.Group("/appointment")
	{
		// Read-only listing/search
		appointment.GET("/list", ar.appointment_handler.ListAppointments)
		appointment.GET("/search", ar.appointment_handler.SearchAppointments)

		// Filters by time range
		appointment.GET("/by-dealership", ar.appointment_handler.ListAppointmentsByDealershipInTimeRange)
		appointment.GET("/by-technician", ar.appointment_handler.ListAppointmentsByTechnicianInTimeRange)
		appointment.GET("/by-bay", ar.appointment_handler.ListAppointmentsByServiceBayInTimeRange)
		appointment.GET("/by-service", ar.appointment_handler.ListAppointmentsByServiceInTimeRange)

		// End-of-day operation
		appointment.POST("/mark-no-show", ar.appointment_handler.MarkNoShowAppointments)
		appointment.PUT("/:id/status", ar.appointment_handler.UpdateAppointmentStatus)
	}
}
