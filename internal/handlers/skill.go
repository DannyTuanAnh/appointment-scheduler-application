package handlers

import (
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/services"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type SkillHandler struct {
	service services.SkillService
}

func NewSkillHandler(service services.SkillService) *SkillHandler {
	return &SkillHandler{service: service}
}

func (h *SkillHandler) ListSkills(ctx *gin.Context) {
	skills, err := h.service.ListSkills(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, skills)
}

func (h *SkillHandler) GetSkillByID(ctx *gin.Context) {
	var req dto.RequestSkillWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	skill, err := h.service.GetSkillByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, skill)
}

func (h *SkillHandler) SearchSkills(ctx *gin.Context) {
	var q dto.SearchSkillsQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	skills, err := h.service.SearchSkillsByName(ctx.Request.Context(), q.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, skills)
}

func (h *SkillHandler) CreateSkill(ctx *gin.Context) {
	var req dto.CreateSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	_, err := h.service.CreateSkill(ctx.Request.Context(), req.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusCreated)
}

func (h *SkillHandler) UpdateSkill(ctx *gin.Context) {
	var reqID dto.RequestSkillWithID
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var req dto.UpdateSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if req.Name == nil {
		utils.ResponseStatusCode(ctx, http.StatusNoContent)
		return
	}

	_, err := h.service.UpdateSkillNameByID(ctx.Request.Context(), reqID.ID, req.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusOK)
}

func (h *SkillHandler) DeleteSkill(ctx *gin.Context) {
	var req dto.RequestSkillWithID
	if err := ctx.ShouldBindUri(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	err := h.service.DeleteSkillByID(ctx.Request.Context(), req.ID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
