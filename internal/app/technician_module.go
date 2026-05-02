package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type TechnicianModule struct {
	routes routes.Routes
}

func NewTechnicianModule() *TechnicianModule {
	// 1. Initialize repository
	repo := repositories.NewTechnicianRepository(db.DB)

	// 2. Initialize service
	service := services.NewTechnicianService(repo)

	// 3. Initialize handler
	handler := handlers.NewTechnicianHandler(service)

	// 4. Initialize routes
	r := routes.NewTechnicianRoutes(handler)

	return &TechnicianModule{routes: r}
}

func (m *TechnicianModule) Routes() routes.Routes {
	return m.routes
}
