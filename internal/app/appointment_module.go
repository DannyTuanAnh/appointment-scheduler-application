package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type AppointmentModule struct {
	routes routes.Routes
}

func NewAppointmentModule() *AppointmentModule {
	// 1. Initialize repository
	appointment_repo := repositories.NewAppointmentRepository(db.DB)

	// 2. Initialize service
	appointment_service := services.NewAppointmentService(appointment_repo)

	// 3. Initialize handler
	appointment_handler := handlers.NewAppointmentHandler(appointment_service)

	// 4. Initialize routes
	appointment_routes := routes.NewAppointmentRoutes(appointment_handler)

	return &AppointmentModule{routes: appointment_routes}
}

func (d *AppointmentModule) Routes() routes.Routes {
	return d.routes
}
