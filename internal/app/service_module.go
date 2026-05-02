package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type ServiceModule struct {
	routes routes.Routes
}

func NewServiceModule() *ServiceModule {
	// 1. Initialize repository
	repo := repositories.NewServiceRepository(db.DB)

	// 2. Initialize service
	service := services.NewServiceService(repo)

	// 3. Initialize handler
	handler := handlers.NewServiceHandler(service)

	// 4. Initialize routes
	r := routes.NewServiceRoutes(handler)

	return &ServiceModule{routes: r}
}

func (m *ServiceModule) Routes() routes.Routes {
	return m.routes
}
