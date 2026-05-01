package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type ServiceBayModule struct {
	routes routes.Routes
}

func NewServiceBayModule() *ServiceBayModule {
	// 1. Initialize repository
	repo := repositories.NewServiceBayRepository(db.DB)

	// 2. Initialize service
	service := services.NewServiceBayService(repo)

	// 3. Initialize handler
	handler := handlers.NewServiceBayHandler(service)

	// 4. Initialize routes
	r := routes.NewServiceBayRoutes(handler)

	return &ServiceBayModule{routes: r}
}

func (m *ServiceBayModule) Routes() routes.Routes {
	return m.routes
}
