package routes

import (
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/handlers"
	"github.com/gin-gonic/gin"
)

type SkillRoutes struct {
	handler *handlers.SkillHandler
}

func NewSkillRoutes(handler *handlers.SkillHandler) Routes {
	return &SkillRoutes{handler: handler}
}

func (sr *SkillRoutes) RegisterApp(r *gin.RouterGroup) {
	skill := r.Group("/skill")
	{
		skill.GET("", sr.handler.ListSkills)
		skill.GET("/search", sr.handler.SearchSkills)
		skill.GET("/:id", sr.handler.GetSkillByID)
	}
}

func (sr *SkillRoutes) RegisterAdmin(r *gin.RouterGroup) {
	skill := r.Group("/skill")
	{
		skill.POST("", sr.handler.CreateSkill)
		skill.PUT("/:id", sr.handler.UpdateSkill)
		skill.DELETE("/:id", sr.handler.DeleteSkill)
	}
}
