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
	skillRepo := repositories.NewSkillRepository(db.DB)
	skillService := services.NewSkillService(skillRepo)
	skillHandler := handlers.NewSkillHandler(skillService)
	skillRoutes := routes.NewSkillRoutes(skillHandler)

	return &SkillModule{routes: skillRoutes}
}

func (sm *SkillModule) Routes() routes.Routes {
	return sm.routes
}
