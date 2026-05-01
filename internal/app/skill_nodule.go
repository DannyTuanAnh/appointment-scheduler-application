package app

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
)

type SkillModule struct {
	routes routes.Routes
}

func NewSkillModule() *SkillModule {
	// 1. Initialize repository
	skillRepo := repositories.NewSkillRepository(db.DB)

	// 2. Initialize service
	skillService := services.NewSkillService(skillRepo)

	// 3. Initialize handler
	skillHandler := handlers.NewSkillHandler(skillService)

	// 4. Initialize routes
	skillRoutes := routes.NewSkillRoutes(skillHandler)

	return &SkillModule{routes: skillRoutes}
}

func (sm *SkillModule) Routes() routes.Routes {
	return sm.routes
}
