package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type DealershipModule struct {
	routes routes.Routes
}

func NewDealershipModule() *DealershipModule {
	// 1. Initialize repository
	dealership_repo := repositories.NewDealershipRepository(db.DB)

	// 2. Initialize service
	dealership_service := services.NewDealershipService(dealership_repo)

	// 3. Initialize handler
	dealership_handler := handlers.NewDealershipHandler(dealership_service)

	// 4. Initialize routes
	dealership_routes := routes.NewDealershipRoutes(dealership_handler)

	return &DealershipModule{routes: dealership_routes}
}

func (d *DealershipModule) Routes() routes.Routes {
	return d.routes
}
